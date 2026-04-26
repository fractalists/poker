import { readFileSync } from "node:fs";
import path from "node:path";

const stylesheet = readFileSync(
  path.resolve(process.cwd(), "src/styles.css"),
  "utf8",
);

describe("table seat state styling", () => {
  it("keeps seat cards compact by tightening chrome without rearranging content", () => {
    expect(stylesheet).toMatch(/\.table-seat\s*\{[^}]*grid-template-rows:\s*auto auto auto;[^}]*gap:\s*8px;[^}]*padding:\s*10px 10px 9px;[^}]*min-height:\s*var\(--orbit-seat-min-height,\s*108px\);/s);
    expect(stylesheet).toMatch(/\.seat-body\s*\{[^}]*gap:\s*6px;[^}]*min-height:\s*0;/s);
    expect(stylesheet).toMatch(/\.seat-footer\s*\{[^}]*gap:\s*12px;[^}]*min-width:\s*0;(?![^}]*margin-top:\s*auto)/s);
    expect(stylesheet).toMatch(/\.seat-stack-pill\s*\{[^}]*min-height:\s*36px;[^}]*padding:\s*6px 9px;/s);
    expect(stylesheet).toMatch(/\.seat-stack-pill--has-result\s*\{[^}]*padding-right:\s*52px;/s);
    expect(stylesheet).not.toMatch(/\.seat-stack-pill--bankroll\s*\{[^}]*padding-right:/s);
    expect(stylesheet).toMatch(/@container table-stage \(max-width:\s*980px\)\s*\{[\s\S]*\.seat-orbit\s*\{[^}]*grid-template-columns:\s*repeat\(auto-fit,\s*minmax\(196px,\s*224px\)\);[^}]*justify-content:\s*center;/s);
    expect(stylesheet).toMatch(/@container table-stage \(max-width:\s*980px\)\s*\{[\s\S]*\.seat-slot \.table-seat\s*\{[^}]*min-height:\s*158px;/s);
  });

  it("lets compact full-ring seat chrome shrink instead of clipping right-side badges", () => {
    expect(stylesheet).toMatch(/\.seat-header\s*\{[^}]*min-width:\s*0;/s);
    expect(stylesheet).toMatch(/\.seat-name-row\s*\{[^}]*min-width:\s*0;/s);
    expect(stylesheet).toMatch(/\.seat-header-meta\s*\{[^}]*min-width:\s*0;/s);
    expect(stylesheet).toMatch(/\.seat-action-pill\s*\{[^}]*max-width:\s*100%;[^}]*overflow:\s*hidden;[^}]*text-overflow:\s*ellipsis;/s);
    expect(stylesheet).toMatch(/\.seat-stack-pill\s*\{[^}]*max-width:\s*100%;[^}]*overflow:\s*hidden;/s);
    expect(stylesheet).toMatch(/\.seat-orbit--8 \.seat-footer,\s*\.seat-orbit--9 \.seat-footer,\s*\.seat-orbit--10 \.seat-footer\s*\{[^}]*display:\s*grid;[^}]*grid-template-columns:\s*minmax\(0,\s*1fr\)\s*max-content;/s);
  });

  it("keeps busted seats in the softer deadened style and makes folded seats slightly more visible", () => {
    expect(stylesheet).toMatch(/\.table-seat\.is-eliminated\s*\{[^}]*border-color:\s*rgba\(150,\s*174,\s*200,\s*0\.08\);[^}]*background:\s*rgba\(4,\s*13,\s*18,\s*0\.18\);[^}]*opacity:\s*0\.42;[^}]*filter:\s*saturate\(0\.24\)\s*brightness\(0\.82\);/s);
    expect(stylesheet).toMatch(/\.table-seat\.is-folded\s*\{[^}]*border-color:\s*rgba\(150,\s*174,\s*200,\s*0\.12\);[^}]*background:\s*linear-gradient\(180deg,\s*rgba\(7,\s*22,\s*27,\s*0\.42\),\s*rgba\(4,\s*16,\s*20,\s*0\.34\)\),[^}]*opacity:\s*0\.64;[^}]*filter:\s*saturate\(0\.46\)\s*brightness\(0\.9\);/s);
  });

  it("uses a quieter folded badge and an even quieter busted badge", () => {
    expect(stylesheet).toMatch(/\.seat-action-pill\.tone-fold\s*\{[^}]*color:\s*rgba\(214,\s*229,\s*219,\s*0\.68\);[^}]*background:\s*rgba\(126,\s*155,\s*139,\s*0\.12\);/s);
    expect(stylesheet).toMatch(/\.seat-action-pill\.tone-out\s*\{[^}]*color:\s*rgba\(182,\s*196,\s*208,\s*0\.56\);[^}]*background:\s*rgba\(112,\s*126,\s*138,\s*0\.08\);/s);
  });

  it("anchors settlement deltas to the lower-right of the stack pill", () => {
    expect(stylesheet).toMatch(/\.seat-stack-pill \.seat-result-pill\s*\{[^}]*right:\s*6px;[^}]*bottom:\s*6px;(?![^}]*top:)[^}]*font-size:\s*0\.66rem;/s);
    expect(stylesheet).toMatch(/\.seat-orbit--8 \.seat-stack-pill--has-result,\s*\.seat-orbit--9 \.seat-stack-pill--has-result,\s*\.seat-orbit--10 \.seat-stack-pill--has-result\s*\{[^}]*padding-right:\s*46px;/s);
    expect(stylesheet).toMatch(/\.seat-orbit--8 \.seat-stack-pill \.seat-result-pill,\s*\.seat-orbit--9 \.seat-stack-pill \.seat-result-pill,\s*\.seat-orbit--10 \.seat-stack-pill \.seat-result-pill\s*\{[^}]*right:\s*5px;[^}]*bottom:\s*5px;(?![^}]*top:)[^}]*font-size:\s*0\.62rem;/s);
  });
});
