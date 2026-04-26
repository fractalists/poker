import { readFileSync } from "node:fs";
import path from "node:path";

const stylesheet = readFileSync(
  path.resolve(process.cwd(), "src/styles.css"),
  "utf8",
);

describe("lobby form styling", () => {
  it("keeps native select menu options readable on the browser popup background", () => {
    expect(stylesheet).toMatch(/\.field select option,\s*select\.field-control option\s*\{[^}]*background:\s*#f8fafc;[^}]*color:\s*#07121f;/s);
    expect(stylesheet).toMatch(/\.field select option:checked,\s*select\.field-control option:checked\s*\{[^}]*background:\s*#93c5fd;[^}]*color:\s*#07121f;/s);
  });
});
