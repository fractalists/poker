import { readFileSync } from "node:fs";
import path from "node:path";

const stylesheet = readFileSync(
  path.resolve(process.cwd(), "src/styles.css"),
  "utf8",
);

describe("room board styling tokens", () => {
  it("renders the board stats as a pure type stack without logos or decorative rails", () => {
    expect(stylesheet).toMatch(/\.board-meta\s*\{[^}]*width:\s*var\(--orbit-stat-width,\s*132px\);[^}]*gap:\s*12px;/s);
    expect(stylesheet).toMatch(/\.table-stat-row\s*\{[^}]*grid-template-columns:\s*1fr;[^}]*padding:\s*0 0 12px;[^}]*border-bottom:\s*1px solid rgba/s);
    expect(stylesheet).toMatch(/\.table-stat-value\s*\{[^}]*font-size:\s*1\.9rem;[^}]*font-weight:\s*600;/s);
  });

  it("anchors the community cards label into the tray as a text tab without an icon", () => {
    expect(stylesheet).toMatch(/\.board-centerpiece\s*\{[^}]*gap:\s*0;/s);
    expect(stylesheet).toMatch(/\.board-badge\s*\{[^}]*margin-bottom:\s*-14px;[^}]*padding:\s*8px 18px 9px;/s);
    expect(stylesheet).toMatch(/\.table-note--board\s*\{[^}]*letter-spacing:\s*0\.18em;[^}]*font-weight:\s*700;/s);
    expect(stylesheet).not.toMatch(/\.board-badge-dot\s*\{/s);
  });

  it("uses a slimmer seat-card size system distinct from the board cards", () => {
    expect(stylesheet).toMatch(/\.seat-card\s*\{[^}]*width:\s*calc\(var\(--seat-card-width,\s*50px\)\s*\*\s*var\(--orbit-card-scale,\s*1\)\);[^}]*height:\s*calc\(var\(--seat-card-height,\s*70px\)\s*\*\s*var\(--orbit-card-scale,\s*1\)\);/s);
    expect(stylesheet).toMatch(/\.seat-card\s+\.card-face\s*\{[^}]*font-size:\s*1\.04rem;[^}]*padding:\s*0 4px;/s);
  });

  it("compacts 8-10 handed seats so full-ring tables can use taller spacing without card overlap", () => {
    expect(stylesheet).toMatch(/\.seat-orbit--8 \.table-seat,\s*\.seat-orbit--9 \.table-seat,\s*\.seat-orbit--10 \.table-seat\s*\{[^}]*padding:\s*10px 10px 10px;[^}]*gap:\s*9px;/s);
    expect(stylesheet).toMatch(/\.seat-orbit--8 \.seat-stack-pill,\s*\.seat-orbit--9 \.seat-stack-pill,\s*\.seat-orbit--10 \.seat-stack-pill\s*\{[^}]*min-height:\s*40px;[^}]*padding:\s*7px 10px;/s);
    expect(stylesheet).toMatch(/\.seat-orbit--8 \.seat-stack-pill \.seat-result-pill,\s*\.seat-orbit--9 \.seat-stack-pill \.seat-result-pill,\s*\.seat-orbit--10 \.seat-stack-pill \.seat-result-pill\s*\{[^}]*top:\s*5px;[^}]*right:\s*5px;[^}]*font-size:\s*0\.62rem;/s);
  });
});
