import { readFileSync } from "node:fs";
import path from "node:path";

const stylesheet = readFileSync(
  path.resolve(process.cwd(), "src/styles.css"),
  "utf8",
);
const cardBackSvg = readFileSync(
  path.resolve(process.cwd(), "src/assets/card-back-diamond.svg"),
  "utf8",
);

describe("card styling tokens", () => {
  it("uses a poker-like card silhouette instead of oversized rounded corners", () => {
    expect(stylesheet).toContain("--card-width: 58px;");
    expect(stylesheet).toContain("--card-height: 82px;");
    expect(stylesheet).toContain("--seat-card-width: 46px;");
    expect(stylesheet).toContain("--seat-card-height: 65px;");
    expect(stylesheet).toContain("--card-radius: 10px;");
  });

  it("uses the selected minimal diamond SVG card back instead of a simple hatch pattern", () => {
    expect(stylesheet).toContain("--card-back-inset: 4px;");
    expect(stylesheet).toContain("--card-back-navy: #102842;");
    expect(stylesheet).toContain("--card-back-gold: #e9bf74;");
    expect(stylesheet).toContain('url("./assets/card-back-diamond.svg")');
    expect(stylesheet).toMatch(/\.card-face--back\s*\{[^}]*padding:\s*var\(--card-back-inset\);/s);
    expect(stylesheet).toMatch(/\.seat-card\s+\.card-face--back\s*\{[^}]*padding:\s*var\(--card-back-inset\);/s);
    expect(stylesheet).toMatch(/\.card-face-pattern\s*\{[^}]*background-color:\s*var\(--card-back-navy\);[^}]*background-image:[^}]*url\("\.\/assets\/card-back-diamond\.svg"\);[^}]*background-repeat:\s*no-repeat;[^}]*background-size:\s*68% auto;/s);
    expect(stylesheet).toMatch(/\.card-face-pattern\s*\{[^}]*box-shadow:[^}]*inset 0 0 0 1px rgba\(231,\s*189,\s*114,\s*0\.78\),[^}]*inset 0 0 0 4px rgba\(8,\s*22,\s*38,\s*0\.92\),[^}]*inset 0 0 0 5px rgba\(231,\s*189,\s*114,\s*0\.48\)/s);
    expect(stylesheet).toMatch(/\.card-face-pattern::before\s*\{[^}]*background:\s*linear-gradient\(115deg,\s*rgba\(255,\s*255,\s*255,\s*0\.14\),\s*transparent 42%\);/s);
    expect(stylesheet).toMatch(/\.card-face-pattern::after\s*\{[^}]*inset:\s*6px;[^}]*border:\s*1px solid rgba\(231,\s*189,\s*114,\s*0\.36\);/s);
    expect(cardBackSvg).toContain('viewBox="0 0 100 146"');
    expect(cardBackSvg).toContain('aria-label="Minimal diamond playing card back"');
    expect(cardBackSvg).toContain("stroke: #e7bd72;");
    expect(cardBackSvg).toContain('id="diamond-frame"');
  });
});
