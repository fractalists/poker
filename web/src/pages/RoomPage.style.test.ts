import { readFileSync } from "node:fs";
import path from "node:path";

const stylesheet = readFileSync(
  path.resolve(process.cwd(), "src/styles.css"),
  "utf8",
);

describe("room board styling tokens", () => {
  it("renders the board stats as a compact table HUD under the community cards", () => {
    expect(stylesheet).toMatch(/\.board-cluster\s*\{[^}]*grid-template-columns:\s*minmax\(0,\s*auto\);[^}]*justify-items:\s*center;/s);
    expect(stylesheet).toMatch(/\.board-centerpiece\s*\{[^}]*width:\s*min\([^}]*var\(--card-width\)\s*\*\s*5[^}]*var\(--board-card-gap,\s*12px\)\s*\*\s*4[^}]*var\(--board-shell-inline-padding,\s*24px\)\s*\*\s*2\)\s*\+\s*2px/s);
    expect(stylesheet).toMatch(/\.board-meta\s*\{[^}]*display:\s*grid;[^}]*grid-template-columns:\s*repeat\(3,\s*max-content\);[^}]*width:\s*fit-content;[^}]*max-width:\s*min\(100%,\s*560px\);[^}]*min-height:\s*36px;[^}]*margin-top:\s*-18px;[^}]*border-radius:\s*999px;/s);
    expect(stylesheet).toMatch(/\.table-stat-row\s*\{[^}]*display:\s*flex;[^}]*align-items:\s*center;[^}]*min-height:\s*34px;[^}]*padding:\s*0 12px;[^}]*border-right:\s*1px solid rgba/s);
    expect(stylesheet).toMatch(/\.table-stat-copy\s*\{[^}]*display:\s*flex;[^}]*align-items:\s*center;/s);
    expect(stylesheet).toMatch(/\.table-stat-label\s*\{[^}]*font-size:\s*0\.5rem;[^}]*letter-spacing:\s*0\.16em;/s);
    expect(stylesheet).toMatch(/\.table-stat-value\s*\{[^}]*font-size:\s*0\.92rem;[^}]*letter-spacing:\s*0;[^}]*font-variant-numeric:\s*tabular-nums lining-nums;/s);
    expect(stylesheet).toMatch(/\.table-stat-row--pot\s+\.table-stat-value\s*\{[^}]*color:\s*var\(--accent\);[^}]*font-size:\s*1\.08rem;/s);
    expect(stylesheet).not.toMatch(/\.board-meta\s*\{[^}]*width:\s*var\(--orbit-stat-width/s);
    expect(stylesheet).not.toMatch(/\.table-stat-row\s*\{[^}]*align-items:\s*baseline;/s);
    expect(stylesheet).not.toMatch(/\.table-stat-row:first-child\s+\.table-stat-value\s*\{/s);
  });

  it("anchors the community cards label into the tray as a text tab without an icon", () => {
    expect(stylesheet).toMatch(/\.board-centerpiece\s*\{[^}]*gap:\s*0;/s);
    expect(stylesheet).toMatch(/\.board-badge\s*\{[^}]*margin-bottom:\s*-18px;[^}]*padding:\s*8px 18px 9px;/s);
    expect(stylesheet).toMatch(/\.table-note--board\s*\{[^}]*letter-spacing:\s*0\.18em;[^}]*font-weight:\s*700;/s);
    expect(stylesheet).not.toMatch(/\.board-badge-dot\s*\{/s);
  });

  it("flips community cards in with a dedicated reveal animation", () => {
    expect(stylesheet).toMatch(/\.community-card\s*\{[^}]*perspective:\s*800px;/s);
    expect(stylesheet).toMatch(/\.community-card\.is-revealing\s+\.card-face\s*\{[^}]*animation:\s*community-card-flip 420ms cubic-bezier\(0\.2,\s*0\.78,\s*0\.24,\s*1\) both;/s);
    expect(stylesheet).toMatch(/@keyframes community-card-flip\s*\{/s);
  });

  it("uses a slimmer seat-card size system distinct from the board cards", () => {
    expect(stylesheet).toMatch(/\.seat-card\s*\{[^}]*width:\s*calc\(var\(--seat-card-width,\s*46px\)\s*\*\s*var\(--orbit-card-scale,\s*1\)\);[^}]*height:\s*calc\(var\(--seat-card-height,\s*65px\)\s*\*\s*var\(--orbit-card-scale,\s*1\)\);/s);
    expect(stylesheet).toMatch(/\.seat-card\s+\.card-face\s*\{[^}]*font-size:\s*1\.04rem;[^}]*padding:\s*0 4px;/s);
  });

  it("compacts 8-10 handed seats so full-ring tables can use taller spacing without card overlap", () => {
    expect(stylesheet).toMatch(/\.seat-orbit--8 \.table-seat,\s*\.seat-orbit--9 \.table-seat,\s*\.seat-orbit--10 \.table-seat\s*\{[^}]*padding:\s*9px;[^}]*gap:\s*8px;/s);
    expect(stylesheet).toMatch(/\.seat-orbit--8 \.seat-stack-pill,\s*\.seat-orbit--9 \.seat-stack-pill,\s*\.seat-orbit--10 \.seat-stack-pill\s*\{[^}]*min-height:\s*34px;[^}]*padding:\s*6px 8px;/s);
    expect(stylesheet).toMatch(/\.seat-orbit--8 \.seat-stack-pill \.seat-result-pill,\s*\.seat-orbit--9 \.seat-stack-pill \.seat-result-pill,\s*\.seat-orbit--10 \.seat-stack-pill \.seat-result-pill\s*\{[^}]*right:\s*5px;[^}]*bottom:\s*5px;(?![^}]*top:)[^}]*font-size:\s*0\.62rem;/s);
  });

  it("stages settlement reveals with dedicated winner and loser animations", () => {
    expect(stylesheet).toMatch(/\.table-seat\.is-settlement-winner\s*\{[^}]*animation:\s*settlement-winner-reveal 1460ms cubic-bezier\(0\.18,\s*0\.82,\s*0\.24,\s*1\) both;/s);
    expect(stylesheet).toMatch(/\.table-seat\.is-settlement-loser\s*\{[^}]*animation:\s*settlement-loser-reveal 1180ms ease-out both;/s);
    expect(stylesheet).toMatch(/\.settlement-panel\.is-animated\s*\{[^}]*animation:\s*settlement-panel-reveal 820ms cubic-bezier\(0\.2,\s*0\.82,\s*0\.24,\s*1\) both;/s);
    expect(stylesheet).toMatch(/\.settlement-entry\s*\{[^}]*animation:\s*settlement-entry-reveal 680ms ease-out both;[^}]*animation-delay:\s*calc\(var\(--settlement-index,\s*0\)\s*\*\s*90ms\s*\+\s*120ms\);/s);
  });

  it("uses table action cue animations for blinds, dealing, and pot collection", () => {
    expect(stylesheet).toMatch(/\.table-action-cue\s*\{[^}]*position:\s*relative;[^}]*z-index:\s*2;[^}]*margin:\s*8px auto 0;[^}]*animation:\s*table-action-cue-in 680ms cubic-bezier\(0\.18,\s*0\.82,\s*0\.24,\s*1\) both;/s);
    expect(stylesheet).toMatch(/\.table-action-cue-label\s*\{[^}]*letter-spacing:\s*0\.2em;[^}]*text-transform:\s*uppercase;/s);
    expect(stylesheet).toMatch(/\.table-action-cue-value\s*\{[^}]*text-overflow:\s*ellipsis;[^}]*text-transform:\s*uppercase;/s);
    expect(stylesheet).not.toMatch(/\.table-action-cue\s*\{[^}]*position:\s*absolute;/s);
    expect(stylesheet).toMatch(/\.table-live-layout--blind-posted\s+\.board-centerpiece::after\s*\{[^}]*animation:\s*table-chip-pulse 860ms ease-out both;/s);
    expect(stylesheet).toMatch(/\.table-live-layout--hole-cards-dealt\s+\.seat-card\s*\{[^}]*animation:\s*seat-card-deal 520ms cubic-bezier\(0\.2,\s*0\.78,\s*0\.24,\s*1\) both;/s);
    expect(stylesheet).toMatch(/\.table-live-layout--pot-collected\s+\.board-cards-shell\s*\{[^}]*animation:\s*pot-collect-flash 780ms ease-out both;/s);
  });

  it("marks disabled controls with a non-interactive cursor", () => {
    expect(stylesheet).toMatch(/button:disabled\s*\{[^}]*cursor:\s*not-allowed;/s);
    expect(stylesheet).toMatch(/\.control-stack button\.start-hand-control\s*\{[^}]*grid-column:\s*1 \/ -1;/s);
  });

  it("switches constrained tables to a playable seat list instead of a clipped orbit", () => {
    expect(stylesheet).toMatch(/\.table-stage\s*\{[^}]*container-type:\s*inline-size;[^}]*container-name:\s*table-stage;/s);
    expect(stylesheet).toMatch(/@container table-stage \(max-width:\s*760px\)\s*\{[\s\S]*?\.seat-orbit\s*\{[^}]*display:\s*grid;[^}]*grid-template-columns:\s*repeat\(auto-fit,\s*minmax\(196px,\s*224px\)\);[^}]*justify-content:\s*center;/s);
    expect(stylesheet).toMatch(/@container table-stage \(max-width:\s*760px\)\s*\{[\s\S]*?\.seat-slot\s*\{[^}]*position:\s*static;[^}]*order:\s*var\(--mobile-seat-order,\s*20\);/s);
  });

  it("keeps desktop seat orbits inside the felt instead of adding an inner vertical scroller", () => {
    expect(stylesheet).toMatch(/\.table-felt\s*\{[^}]*height:\s*100%;[^}]*overflow:\s*hidden;/s);
    expect(stylesheet).toMatch(/\.table-live-layout\s*\{[^}]*height:\s*min\(100%,\s*var\(--orbit-min-height,\s*560px\)\);[^}]*min-height:\s*0;/s);
    expect(stylesheet).toMatch(/\.seat-orbit\s*\{[^}]*height:\s*100%;[^}]*min-height:\s*0;/s);
    expect(stylesheet).toMatch(/@container table-stage \(max-width:\s*760px\)\s*\{[\s\S]*?\.table-felt\s*\{[^}]*overflow:\s*auto;/s);
    expect(stylesheet).toMatch(/@container table-stage \(max-width:\s*760px\)\s*\{[\s\S]*?\.table-live-layout\s*\{[^}]*height:\s*auto;[^}]*min-height:\s*auto;/s);
  });

  it("uses a mobile control layout with a fixed action dock for portrait play", () => {
    expect(stylesheet).toMatch(/@media \(max-width:\s*900px\),\s*\(orientation:\s*portrait\)\s*\{/s);
    expect(stylesheet).toMatch(/\.room-shell\.has-pending-action\s*\{[^}]*padding-bottom:\s*min\(50svh,\s*430px\);/s);
    expect(stylesheet).toMatch(/\.room-shell\.has-pending-action\s+\.action-bar\s*\{[^}]*position:\s*fixed;[^}]*bottom:\s*12px;[^}]*max-height:\s*min\(50svh,\s*430px\);/s);
    expect(stylesheet).toMatch(/\.room-shell\.has-pending-action\s+\.room-grid\s+\.room-history-stack--scroll\s*\{[^}]*display:\s*none;/s);
    expect(stylesheet).toMatch(/\.room-shell\.has-pending-action\s+\.control-stack\s*\{[^}]*display:\s*none;/s);
    expect(stylesheet).toMatch(/\.action-buttons\s*\{[^}]*grid-template-columns:\s*repeat\(2,\s*minmax\(0,\s*1fr\)\);/s);
    expect(stylesheet).toMatch(/\.table-stage\s*\{[^}]*min-height:\s*240px;/s);
    expect(stylesheet).toMatch(/\.table-felt\s*\{[^}]*min-height:\s*222px;/s);
    expect(stylesheet).toMatch(/\.table-action-cue\s*\{[^}]*max-width:\s*calc\(100%\s*-\s*24px\);[^}]*margin:\s*8px auto 0;/s);
    expect(stylesheet).toMatch(/\.board-cards-shell\s*\{[^}]*padding:\s*24px var\(--board-shell-inline-padding\);[^}]*border-radius:\s*22px;/s);
  });

  it("uses a portrait-first seat rail instead of stacking every player as full-width rows", () => {
    expect(stylesheet).toMatch(/@media \(max-width:\s*900px\),\s*\(orientation:\s*portrait\)\s*\{[\s\S]*?\.seat-orbit\s*\{[^}]*display:\s*flex;[^}]*overflow-x:\s*auto;[^}]*scroll-snap-type:\s*x mandatory;/s);
    expect(stylesheet).toMatch(/@media \(max-width:\s*900px\),\s*\(orientation:\s*portrait\)\s*\{[\s\S]*?\.seat-slot\s*\{[^}]*flex:\s*0 0 clamp\(174px,\s*58vw,\s*220px\);[^}]*scroll-snap-align:\s*start;/s);
    expect(stylesheet).not.toMatch(/@media \(max-width:\s*900px\),\s*\(orientation:\s*portrait\)\s*\{[\s\S]*?\.seat-orbit\s*\{[^}]*grid-template-columns:\s*1fr;/s);
  });
});
