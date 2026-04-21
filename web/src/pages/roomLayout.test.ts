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

    expect(px(layout.spec.seatWidth)).toBeGreaterThanOrEqual(188);
    expect(px(layout.spec.seatMinHeight)).toBeGreaterThanOrEqual(170);
    expect(layout.spec.cardScale).toBeGreaterThanOrEqual(0.94);
    expect(px(layout.spec.minHeight)).toBeGreaterThanOrEqual(872);
  });
});
