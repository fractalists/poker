import { readFileSync } from "node:fs";
import path from "node:path";

const stylesheet = readFileSync(
  path.resolve(process.cwd(), "src/styles.css"),
  "utf8",
);

describe("card styling tokens", () => {
  it("uses a poker-like card silhouette instead of oversized rounded corners", () => {
    expect(stylesheet).toContain("--card-width: 58px;");
    expect(stylesheet).toContain("--card-height: 82px;");
    expect(stylesheet).toContain("--card-radius: 10px;");
  });

  it("keeps the card back frame tight and the pattern scaled like a poker back", () => {
    expect(stylesheet).toContain("--card-back-inset: 4px;");
    expect(stylesheet).toContain("--card-back-pattern-size: 10px;");
    expect(stylesheet).toMatch(/\.card-face--back\s*\{[^}]*padding:\s*var\(--card-back-inset\);/s);
    expect(stylesheet).toMatch(/\.card-face-pattern\s*\{[^}]*background-size:\s*var\(--card-back-pattern-size\)\s+var\(--card-back-pattern-size\);/s);
    expect(stylesheet).toMatch(/\.card-face-pattern\s*\{[^}]*background-position:\s*center;[^}]*background-image:\s*repeating-linear-gradient/s);
    expect(stylesheet).toMatch(/\.card-face-pattern::after\s*\{[^}]*border:\s*1px solid rgba\(255,\s*255,\s*255,\s*0\.18\);/s);
  });
});
