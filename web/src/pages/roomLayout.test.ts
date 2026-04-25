import { buildOrbitLayout } from "./roomLayout";

function px(value: string) {
  return Number.parseFloat(value.replace("px", ""));
}

describe("roomLayout", () => {
  it("keeps seat anchors inside the computed safe inset for 2, 6, and 10 handed tables", () => {
    [2, 6, 10].forEach((playerCount) => {
      const layout = buildOrbitLayout(playerCount, playerCount - 1);

      expect(layout.positions.size).toBe(playerCount);

      for (const position of layout.positions.values()) {
        expect(position.x).toBeGreaterThanOrEqual(layout.safeInset.left);
        expect(position.x).toBeLessThanOrEqual(100 - layout.safeInset.right);
        expect(position.y).toBeGreaterThanOrEqual(layout.safeInset.top);
        expect(position.y).toBeLessThanOrEqual(100 - layout.safeInset.bottom);
      }
    });
  });

  it("keeps full-ring seats above the readability floor", () => {
    const layout = buildOrbitLayout(10, 9);

    expect(px(layout.spec.seatWidth)).toBeGreaterThanOrEqual(184);
    expect(px(layout.spec.seatMinHeight)).toBeGreaterThanOrEqual(152);
    expect(layout.spec.cardScale).toBeGreaterThanOrEqual(0.94);
    expect(px(layout.spec.minHeight)).toBeGreaterThanOrEqual(804);
  });

  it("packs six-handed tables into a denser but still readable seat footprint", () => {
    const layout = buildOrbitLayout(6, 5);

    expect(px(layout.spec.seatWidth)).toBeLessThanOrEqual(200);
    expect(px(layout.spec.seatMinHeight)).toBeLessThanOrEqual(158);
    expect(px(layout.spec.minHeight)).toBeLessThanOrEqual(684);
  });
});
