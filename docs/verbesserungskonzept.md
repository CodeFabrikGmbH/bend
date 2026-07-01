# Verbesserungskonzept „bend"

> Stand: 2026-07-01 · Status: Entwurf zur Abstimmung

## Zweck dieses Dokuments

`bend` ist funktional, wurde aber lange nicht weiterentwickelt. Dieses Dokument
fasst die festgestellten Probleme in **Sicherheit**, **Bedienung/Performance** und
**technischem Reifegrad** zusammen und beschreibt ein priorisiertes Vorgehen zur
Verbesserung.

Bewusst akzeptierte Entscheidungen (kein Handlungsbedarf) sind am Ende unter
[Bewusst akzeptiert](#bewusst-akzeptiert) dokumentiert, damit sie nicht immer
wieder neu diskutiert werden.

---

## Was die Anwendung ist

`bend` ist ein **Request-Catcher / Mock-Proxy** in Go: Eingehende Requests auf
beliebige Pfade werden gespeichert (BoltDB), im Dashboard sichtbar gemacht und
können später an ein Ziel weitergeleitet werden. Über `/configs` lassen sich pro
Pfad Mock-Antworten oder Forward-Ziele (inkl. Regex-Pfade) definieren.
Authentifizierung optional über Keycloak.

Die Schichtenarchitektur (`domain` / `application` / `infrastructure`) ist sauber
getrennt und eine gute Grundlage. Die Probleme liegen nicht in der Struktur,
sondern in Sicherheit, Bedienung und Reifegrad.

---

## 1. Sicherheit

| # | Befund | Ort | Wirkung |
|---|--------|-----|---------|
| **S1** | **Login sendet Benutzername + Passwort als GET-Query-Parameter.** Das Formular hat kein `method="post"`, der Handler liest `query.Get("username"/"password")`. | `resources/login.html:29,38`, `infrastructure/httpHandler/loginPage.go:46-47` | Klartext-Credentials landen in URL, Browser-History, Server-Logs und Referrer. |
| **S2** | **`/api/requests/` hat keine Authentifizierung.** `RequestAPI` besitzt kein `KeyCloakService` (im Gegensatz zu `ConfigAPI`). | `main.go:51`, `infrastructure/httpHandler/requestAPI.go` | Jeder kann gespeicherte Requests löschen **und** Requests an eine beliebige URL weiterleiten lassen → SSRF / offener Proxy. |
| **S6** | **Kein Body-Size-Limit beim Tracking**; der gesamte Body wird per `ReadAll` in den Speicher gelesen. | `infrastructure/httpHandler/trackRequest.go:27` | Große/böswillige Uploads → OOM (DoS). |
| **S7** | **Session-Cookies ohne Schutzflags** (`HttpOnly`, `Secure`, `SameSite`), Ablauf 365 Tage. | `infrastructure/jwt/keycloak/service.go:117-121` | Token per JS auslesbar (bei XSS), über HTTP übertragbar, keine echte Session-Begrenzung. |

**Maßnahmen**
- **S1:** Formular auf `method="post"` umstellen, Handler liest `r.PostFormValue(...)`.
- **S2:** `/api/requests/` hinter dieselbe Authentifizierung legen wie `ConfigAPI`.
- **S6:** `http.MaxBytesReader` für eingehende Bodies + Server-Timeouts
  (`ReadTimeout`/`WriteTimeout`) setzen.
- **S7:** Cookies mit `HttpOnly`, `Secure`, `SameSite=Lax`; Laufzeit an
  Token-Expiry koppeln.

---

## 2. Bedienung & Performance

### 2.1 Performance: Dashboard wird mit wachsender Datenmenge langsam ⚠️

Das gravierendste Bedienungsproblem: Das Dashboard wird mit der Zeit spürbar
langsam, weil bei **jedem** Seitenaufruf zu viele Daten geladen werden.

**Ursache 1 — Pfad-Übersicht zählt die gesamte DB durch.**
`getPaths()` holt alle Pfade und ruft pro Pfad `GetRequestCountForPath()` auf.
Diese Methode zählt via `ForEach` **jeden einzelnen Eintrag** einzeln durch
(`infrastructure/boltDB/requestRepository.go:94-112`). Ergebnis: N Pfade × M
Requests = **Full-Scan der gesamten Datenbank bei jedem Dashboard-Aufruf**
(`application/dashboardService.go:77-95`).

**Ursache 2 — Request-Auswahl deserialisiert alle vollständigen Requests.**
`GetRequestsForPath()` lädt und deserialisiert **jeden vollständigen Request**
eines Pfads inkl. Body, aller Header und der kompletten Response
(`infrastructure/boltDB/requestRepository.go:72-92`) — obwohl im Dropdown nur
`ID` und `Timestamp` angezeigt werden (`RequestAbstract`,
`application/dashboardService.go:14-17`). Der komplette Rest wird sofort verworfen.

**Ursache 3 — keine Paginierung.** Es gibt keine Begrenzung; die Datenmenge
wächst monoton, jeder Aufruf lädt alles.

**Maßnahmen**
- **Zählen ohne Full-Scan:** Anzahl über `bucket.Stats().KeyN` bzw. einen
  gepflegten Zähler bestimmen statt jeden Eintrag einzeln zu iterieren.
- **Leichtgewichtige Liste:** Für die Request-Liste nur eine schlanke Projektion
  laden (ID, Timestamp, Methode, Statuscode, Größe) statt der vollen Objekte.
  Optionen: separater Index/Summary-Eintrag beim Speichern, oder nur die für die
  Liste nötigen Felder deserialisieren. Vollständige Details erst beim Öffnen
  eines einzelnen Requests laden (`GetRequest`).
- **Paginierung:** BoltDB-Cursor (`Seek`/Range) nutzen und nur eine Seite laden;
  Pfad-Counts bei vielen Pfaden ggf. cachen oder lazy nachladen.

### 2.2 Weitere Bedienungsprobleme

- **Keine Navigation.** Kein Menü, keine Verlinkung zwischen Dashboard ↔ Configs
  ↔ Readme. Man muss URLs kennen; einziges UI-Element ist ein „?"-Button.
- **Alles über `alert()` / `confirm()`.** Das Ergebnis von „Send Request" wird als
  roher JSON-String in eine Alert-Box geworfen
  (`resources/dashboard.html:68`). Kein Loading-State, keine Inline-Fehler.
- **Request-Auswahl nur als `<select>`-Dropdown** mit Timestamp
  (`resources/dashboard.html:138`). Methode/Statuscode nicht auf einen Blick
  sichtbar; bei vielen Requests unbrauchbar.
- **Kein Live-Update.** Für einen Request-Catcher ist „Live-Tail" die
  Kernfunktion — hier muss man manuell neu laden, um neue Requests zu sehen.
- **Kein Suchen / Filtern.** Pfad-Tabelle und Dropdown skalieren nicht.
- **Body & Header roh in Tabellenzellen** — kein JSON-Pretty-Print, kein
  Copy-Button, kein Umbruch-Handling; lange Bodys sprengen das Layout.
- **Nicht responsive.** `display: table-cell`-Layout, feste Breiten
  (`500px`, `200px`) → auf schmalem Fenster/Mobil kaputt.
- **Login-Seite ist ungestylt.** Sie referenziert Klassen (`.theme-button`,
  `.content-group`, `.label-column`), die nirgends definiert sind, und bindet
  `styles.css` nicht ein.
- **Accessibility & Kleinkram:** Klick-Handler auf `<td>` (nicht per Tastatur
  bedienbar), keine ARIA-Labels, leerer `<title>` auf allen Seiten,
  Regex-Pfade ohne Test-/Vorschau-Möglichkeit, kein Hinweis, welche Config
  gematcht hat.

### 2.3 Zielbild: Der Live-Inspector

Die heutige Oberfläche ist **seitenbasiert** — Pfad-Tabelle → Timestamp-Dropdown
→ Detail-Tabelle, jeder Schritt ein vollständiger Reload. Das Produkt ist aber
eigentlich ein **Live-Inspector**: Webhook draufzeigen, zusehen, wie Requests
eintreffen, reinklicken, weiterleiten. Das neue Konzept bildet genau diesen
Workflow ab.

> Interaktiver Mockup: `docs/bend-mockup.html` (bzw. veröffentlichter Artifact).
> Er ist funktional — Endpoints wählen, Requests öffnen, Tabs, Live-Streaming
> mit Toasts, Regex-Tester.

**Leitidee.** Weg vom Seiten-Modell, hin zu **einem Screen mit drei Spalten**
(wie Linear / Postman / webhook.site):

```
┌──────────── Topbar: bend · Requests|Endpoints|Docs · Catch-URL [copy] · User ┐
├───────────────┬──────────────────────────┬──────────────────────────────────┤
│  ENDPOINTS    │  REQUESTS (Live-Liste)    │  DETAIL                           │
│  /webhooks/…  │  POST  /webhooks/…  200   │  POST /webhooks/stripe           │
│  /gh/push     │  POST  /webhooks/…  500   │  [An Ziel senden][Replay][cURL]  │
│  /api/aggr…•  │  GET   /webhooks/…  200   │  ─ Request | Response ─           │
│  + Endpoint   │  ⟳ live · Filter…         │  { pretty-printed body }         │
└───────────────┴──────────────────────────┴──────────────────────────────────┘
```

- **Spalte 1 — Endpoints:** die „Kanäle" mit Request-Count und Aktiv-Indikator;
  Regex-Endpoints sind als solche markiert.
- **Spalte 2 — Request-Liste:** chronologisch, neueste oben. Jede Zeile zeigt
  **Methoden-Chip (farbcodiert), Status-Chip, Zeit, Größe, User-Agent** auf einen
  Blick. Filterfeld + Live-Toggle.
- **Spalte 3 — Detail:** Kopf mit Methode/Pfad/Status + Aktions-Toolbar; Tabs
  **Request / Response**; Body **pretty-printed und syntaxgefärbt**, Header als
  Tabelle, Copy-Buttons.

**Konkrete Prinzipien**

- **Navigation:** persistente Topbar (Requests / Endpoints / Docs) und die
  **Catch-URL prominent zum Kopieren** — das erste, was man zum Loslegen braucht.
  Gesetzter `<title>` je Ansicht.
- **Live-Tail** via Server-Sent Events (SSE) — kleinste Erweiterung, kein
  WebSocket nötig. Neue Requests erscheinen automatisch (mit Toggle zum
  Pausieren), keine manuellen Reloads mehr.
- **Feedback ohne `alert`/`confirm`:** Aktionen (An Ziel senden, Replay,
  cURL kopieren, Löschen) geben **Inline-Toasts** mit Loading-/Erfolg-/Fehler-State;
  `fetch` statt `XMLHttpRequest`.
- **Fehler sichtbar machen:** Forwarding-Fehler (z.&nbsp;B. `connection refused`)
  und 4xx/5xx werden im Response-Tab klar rot ausgewiesen statt in einer
  Alert-Box zu verschwinden.
- **Endpoint-Editor mit Live-Regex-Tester:** Test-Pfad eingeben → sofort
  Match/kein-Match + resultierende Ziel-URL. Mock-Antwort vs. Weiterleitung sind
  klar getrennt.
- **Responsiv & barrierearm:** Drei-Spalten-Grid, das auf schmalen Screens zu
  einem navigierbaren Stack zusammenfällt; echte Buttons statt klickbarer `<td>`,
  Fokus-States, `prefers-reduced-motion` respektiert. Login-Seite bekommt
  Styling.

**Visuelle Sprache (bewusste Entscheidungen)**

- **Dark-first + Monospace für alle Daten** (Pfade, Header, Bodies) — die
  Vernakular­sprache von Dev-Tools; authentischer als eine dekorative Schrift.
  System-UI-Sans fürs Chrome, System-Mono für Daten (kein Font-Embedding nötig).
- **Semantische Farben getrennt vom Akzent:** HTTP-Methoden (GET grün, POST blau,
  PUT amber, PATCH violett, DELETE rot) und Statusklassen (2xx/4xx/5xx) sind
  Daten-Kodierung — man liest Zustand, bevor man Zahlen liest.
- **Kontinuität:** der Teal-Akzent (`#34c3d6`) führt die vorhandene
  Header-Gradient-Note verfeinert weiter, statt sie zu verwerfen.

**Wie das Zielbild die Probleme aus 2.1 & 2.2 löst**

| Problem | Lösung im Konzept |
|---------|-------------------|
| Dashboard wird langsam (2.1) | Liste lädt nur schlanke Projektion + Paginierung; Vollobjekt erst beim Öffnen eines Requests |
| Keine Navigation | Persistente Topbar + kopierbare Catch-URL |
| Timestamp-Dropdown | Scannbare Live-Liste mit Methode/Status/Zeit/Größe |
| Kein Live-Update | Live-Tail (SSE) mit Auto-Einblendung |
| `alert`/roher JSON | Toasts + pretty-printed, gefärbte Bodies mit Copy |
| Regex blind | Live-Regex-Tester im Editor |
| Nicht responsive / ungestylt | Dark-Theme, kollabierendes Grid, gestylte Login-Seite |

---

## 3. Technischer Reifegrad

- **Go 1.15** und veraltete/unmaintainte Abhängigkeiten: `boltdb/bolt`
  (→ `etcd-io/bbolt`), `gocloak/v7`, `dgrijalva/jwt-go` (archiviert, CVE-behaftet).
  Überall veraltetes `ioutil`.
- **Templates werden bei jedem Request von der Platte gelesen und neu geparst**
  (`template.Must(template.ParseFiles(...))`,
  `infrastructure/htmlTemplate/template.go:9`). `Must` würde bei Parse-Fehler
  paniken. Sollte einmalig via `embed.FS` eingebettet und vorkompiliert werden.
- **Routing per String-Parsing** (`TrimPrefix`, `LastIndex("/")`,
  `dashboardPage.go:30`, `requestAPI.go:60`). Fragil bei verschachtelten Pfaden.
  Ein Router (chi/gorilla) mit Pfad-Parametern wäre robuster.
- **Fehler werden verschluckt** (`_ =` an Save/Marshal/Write) und **`recover()`
  als Catch-all** in jedem Handler, das nur `fmt.Println` macht → keine echten
  Fehlerpfade, kein strukturiertes Logging.
- **`FindMatchingConfig` kompiliert bei jedem eingehenden Request für jede Config
  die Regex neu** (`domain/config/config.go:42`) und macht einen vollen O(n)-Scan.
  Regex sollten beim Speichern kompiliert/gecacht werden.
- **HTTP-Server ohne Timeouts** (`ListenAndServe`) → Slowloris-anfällig; kein
  Graceful Shutdown; Port hart auf `8080`.
- **Ad-hoc-Migration** in `main.go:61` statt eines Migrationskonzepts.

**Maßnahmen**
- Go aktualisieren, `bbolt` + gepflegte Deps, `io` statt `ioutil`,
  `dgrijalva/jwt-go` ablösen.
- Templates und statische Assets via `embed.FS` einbetten und einmalig parsen
  (Single-Binary-Deployment).
- Router mit Pfad-Parametern einführen; String-Parsing entfernen.
- Strukturiertes Logging (`slog`), echte Fehlerbehandlung statt `recover`-Catch-all;
  Graceful Shutdown; Port/Config über Env.
- Regex-Configs beim Speichern kompilieren/cachen.

---

## 4. Roadmap

| Phase | Inhalt | Aufwand | Nutzen |
|-------|--------|---------|--------|
| **0 — Sofort** | S1, S2, S6, S7 (Login-POST, API-Auth, Body-Limit + Timeouts, Cookie-Flags) | klein | Schließt kritische Sicherheitslücken |
| **1 — Fundament** | embed.FS + Template-Caching, Deps/Go-Update, Logging, Router | mittel | Stabil, wartbar, Single-Binary |
| **2 — Performance & UX-Kern** | Dashboard-Performance (2.1), Layout+Navigation, Request-Liste, fetch/Toasts, JSON-Rendering | mittel | Behebt die Verlangsamung + spürbarer Bedienungssprung |
| **3 — Komfort** | Live-Tail (SSE), Suche/Filter, Regex-Tester, Responsive | mittel | Angenehm & konkurrenzfähig |

Phase 0 ist unabhängig vom Rest und sollte zuerst kommen — kleine, klar
abgegrenzte Änderungen mit hohem Risiko-Payoff. Das Performance-Problem (2.1) ist
das wichtigste Bedienungsthema und in Phase 2 eingeplant; es kann bei Bedarf
vorgezogen werden, da es unabhängig vom UI-Redesign umsetzbar ist.

---

## Bewusst akzeptiert

Diese Punkte wurden geprüft und bewusst so belassen — **kein Handlungsbedarf**:

- **TLS-Verifikation beim Forwarding deaktiviert** (`InsecureSkipVerify: true`,
  `infrastructure/http/transport.go:113`). Weiterleitungen gehen typischerweise an
  Entwickler-Rechner mit self-signed HTTPS-Zertifikaten; das Flag ist dafür der
  pragmatische Weg. (Hinweis: Auf reine `http://`-Ziele hat das Flag ohnehin keine
  Wirkung — es greift nur beim TLS-Handshake von `https://`-Zielen.)
- **Authentifizierung komplett aus, wenn `KEYCLOAK_HOST` leer**
  (`infrastructure/jwt/keycloak/service.go:41`). Bewusste Entscheidung für
  lokalen/Dev-Betrieb.
