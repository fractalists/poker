import { render, screen } from "@testing-library/react";

import { CardFace } from "./CardFace";

describe("CardFace", () => {
  it("renders hidden cards as a card back instead of literal mask text", () => {
    const { container } = render(<CardFace card="**" />);

    expect(screen.queryByText("**")).not.toBeInTheDocument();
    expect(container.querySelector(".card-face--back")).toBeInTheDocument();
    expect(container.querySelector(".card-face-pattern")).toBeInTheDocument();
  });
});
