import { fireEvent, render, screen } from "@testing-library/react";
import { vi } from "vitest";

import { ActionBar } from "./ActionBar";

describe("ActionBar", () => {
  afterEach(() => {
    vi.useRealTimers();
  });

  it("renders numeric action labels directly on the primary buttons", () => {
    render(
      <ActionBar
        roomId="room-001"
        pot={8}
        pendingAction={{
          token: "turn-1",
          seatIndex: 5,
          minAmount: 2,
          minBetAmount: 4,
          maxAmount: 96,
          canCheck: false,
          canCall: true,
          canBet: true,
          canFold: true,
          canAllIn: true,
        }}
        onSubmit={vi.fn()}
      />,
    );

    expect(screen.getByRole("button", { name: "Fold" })).toBeInTheDocument();
    expect(screen.getByRole("button", { name: "Call 2" })).toBeInTheDocument();
    expect(screen.getByRole("button", { name: "Raise to 4~96" })).toBeInTheDocument();
    expect(screen.getByRole("button", { name: "All-in 96" })).toBeInTheDocument();
    expect(screen.queryByText(/seat 6/i)).not.toBeInTheDocument();
    expect(screen.queryByText("To call 2")).not.toBeInTheDocument();
    expect(screen.queryByText("Min raise to 4")).not.toBeInTheDocument();
    expect(screen.queryByText("Stack 96")).not.toBeInTheDocument();
    expect(screen.queryByLabelText(/bet amount/i)).not.toBeInTheDocument();
  });

  it("expands betting controls instead of submitting immediately when bet is clicked", () => {
    const onSubmit = vi.fn();

    render(
      <ActionBar
        roomId="room-001"
        pot={16}
        pendingAction={{
          token: "turn-1",
          seatIndex: 5,
          minAmount: 2,
          minBetAmount: 4,
          maxAmount: 96,
          canCheck: false,
          canCall: true,
          canBet: true,
          canFold: true,
          canAllIn: true,
        }}
        onSubmit={onSubmit}
      />,
    );

    fireEvent.click(screen.getByRole("button", { name: "Raise to 4~96" }));

    expect(onSubmit).not.toHaveBeenCalled();
    expect(screen.getByLabelText(/raise to amount/i)).toBeInTheDocument();
    expect(screen.getByRole("button", { name: "1/4 Pot" })).toBeInTheDocument();
    expect(screen.getByRole("button", { name: "1/2 Pot" })).toBeInTheDocument();
    expect(screen.getByRole("button", { name: "1 Pot" })).toBeInTheDocument();
    expect(screen.getByRole("button", { name: "Confirm raise" })).toBeInTheDocument();
  });

  it("uses quick bet sizing and confirms the selected amount", () => {
    const onSubmit = vi.fn();

    render(
      <ActionBar
        roomId="room-001"
        pot={24}
        pendingAction={{
          token: "turn-1",
          seatIndex: 5,
          minAmount: 2,
          minBetAmount: 4,
          maxAmount: 96,
          canCheck: false,
          canCall: true,
          canBet: true,
          canFold: true,
          canAllIn: true,
        }}
        onSubmit={onSubmit}
      />,
    );

    fireEvent.click(screen.getByRole("button", { name: "Raise to 4~96" }));
    fireEvent.click(screen.getByRole("button", { name: "1/2 Pot" }));
    fireEvent.click(screen.getByRole("button", { name: "Confirm raise" }));

    expect(onSubmit).toHaveBeenCalledWith({ token: "turn-1", actionType: "BET", amount: 15 });
  });

  it("updates the bet input with ceiling-rounded shortcut sizes", () => {
    render(
      <ActionBar
        roomId="room-001"
        pot={9}
        pendingAction={{
          token: "turn-1",
          seatIndex: 5,
          minAmount: 0,
          minBetAmount: 1,
          maxAmount: 96,
          canCheck: true,
          canCall: false,
          canBet: true,
          canFold: true,
          canAllIn: true,
        }}
        onSubmit={vi.fn()}
      />,
    );

    fireEvent.click(screen.getByRole("button", { name: "Bet 1~96" }));
    expect((screen.getByLabelText(/bet amount/i) as HTMLInputElement).value).toBe("");
    fireEvent.click(screen.getByRole("button", { name: "1/4 Pot" }));
    expect(screen.getByLabelText(/bet amount/i)).toHaveValue(3);

    fireEvent.click(screen.getByRole("button", { name: "1/2 Pot" }));
    expect(screen.getByLabelText(/bet amount/i)).toHaveValue(5);
  });

  it("uses call-aware pot shortcuts when there is already a live amount to call", () => {
    render(
      <ActionBar
        roomId="room-001"
        pot={3}
        pendingAction={{
          token: "turn-1",
          seatIndex: 5,
          minAmount: 2,
          minBetAmount: 4,
          maxAmount: 100,
          canCheck: false,
          canCall: true,
          canBet: true,
          canFold: true,
          canAllIn: true,
        }}
        onSubmit={vi.fn()}
      />,
    );

    fireEvent.click(screen.getByRole("button", { name: "Raise to 4~100" }));

    fireEvent.click(screen.getByRole("button", { name: "1/4 Pot" }));
    expect(screen.getByLabelText(/raise to amount/i)).toHaveValue(4);

    fireEvent.click(screen.getByRole("button", { name: "1/2 Pot" }));
    expect(screen.getByLabelText(/raise to amount/i)).toHaveValue(5);

    fireEvent.click(screen.getByRole("button", { name: "1 Pot" }));
    expect(screen.getByLabelText(/raise to amount/i)).toHaveValue(7);
  });

  it("keeps shortcut selection responsive even when two ratios clamp to the same minimum bet", () => {
    render(
      <ActionBar
        roomId="room-001"
        pot={2}
        pendingAction={{
          token: "turn-1",
          seatIndex: 5,
          minAmount: 2,
          minBetAmount: 4,
          maxAmount: 100,
          canCheck: false,
          canCall: true,
          canBet: true,
          canFold: true,
          canAllIn: true,
        }}
        onSubmit={vi.fn()}
      />,
    );

    fireEvent.click(screen.getByRole("button", { name: "Raise to 4~100" }));
    fireEvent.click(screen.getByRole("button", { name: "1/2 Pot" }));

    expect(screen.getByLabelText(/raise to amount/i)).toHaveValue(4);
    expect(screen.getByRole("button", { name: "1/2 Pot" })).toHaveAttribute("aria-pressed", "true");
    expect(screen.getByRole("button", { name: "1/4 Pot" })).toHaveAttribute("aria-pressed", "false");

    fireEvent.click(screen.getByRole("button", { name: "1/4 Pot" }));

    expect(screen.getByLabelText(/raise to amount/i)).toHaveValue(4);
    expect(screen.getByRole("button", { name: "1/4 Pot" })).toHaveAttribute("aria-pressed", "true");
    expect(screen.getByRole("button", { name: "1/2 Pot" })).toHaveAttribute("aria-pressed", "false");
  });

  it("keeps the chosen shortcut amount when the same pending action rerenders", () => {
    const { rerender } = render(
      <ActionBar
        roomId="room-001"
        pot={18}
        pendingAction={{
          token: "turn-1",
          seatIndex: 5,
          minAmount: 2,
          minBetAmount: 6,
          maxAmount: 98,
          canCheck: false,
          canCall: true,
          canBet: true,
          canFold: true,
          canAllIn: true,
        }}
        onSubmit={vi.fn()}
      />,
    );

    fireEvent.click(screen.getByRole("button", { name: "Raise to 6~98" }));
    fireEvent.click(screen.getByRole("button", { name: "1/2 Pot" }));
    expect(screen.getByLabelText(/raise to amount/i)).toHaveValue(12);

    rerender(
      <ActionBar
        roomId="room-001"
        pot={18}
        pendingAction={{
          token: "turn-1",
          seatIndex: 5,
          minAmount: 2,
          minBetAmount: 6,
          maxAmount: 98,
          canCheck: false,
          canCall: true,
          canBet: true,
          canFold: true,
          canAllIn: true,
        }}
        onSubmit={vi.fn()}
      />,
    );

    expect(screen.getByLabelText(/raise to amount/i)).toHaveValue(12);
  });

  it("lets the bet input stay empty while the user edits it", () => {
    render(
      <ActionBar
        roomId="room-001"
        pot={24}
        pendingAction={{
          token: "turn-1",
          seatIndex: 5,
          minAmount: 2,
          minBetAmount: 4,
          maxAmount: 96,
          canCheck: false,
          canCall: true,
          canBet: true,
          canFold: true,
          canAllIn: true,
        }}
        onSubmit={vi.fn()}
      />,
    );

    fireEvent.click(screen.getByRole("button", { name: "Raise to 4~96" }));
    fireEvent.change(screen.getByLabelText(/raise to amount/i), { target: { value: "" } });

    expect((screen.getByLabelText(/raise to amount/i) as HTMLInputElement).value).toBe("");
  });

  it("fills the input with the full pot size when no call is pending", () => {
    render(
      <ActionBar
        roomId="room-001"
        pot={24}
        pendingAction={{
          token: "turn-1",
          seatIndex: 5,
          minAmount: 0,
          minBetAmount: 1,
          maxAmount: 96,
          canCheck: true,
          canCall: false,
          canBet: true,
          canFold: true,
          canAllIn: true,
        }}
        onSubmit={vi.fn()}
      />,
    );

    fireEvent.click(screen.getByRole("button", { name: "Bet 1~96" }));
    fireEvent.click(screen.getByRole("button", { name: "1 Pot" }));

    expect(screen.getByLabelText(/bet amount/i)).toHaveValue(24);
  });

  it("hides call when the required amount is higher than the remaining stack", () => {
    render(
      <ActionBar
        roomId="room-001"
        pot={119}
        pendingAction={{
          token: "turn-1",
          seatIndex: 5,
          minAmount: 85,
          minBetAmount: 0,
          maxAmount: 76,
          canCheck: false,
          canCall: true,
          canBet: false,
          canFold: true,
          canAllIn: true,
        }}
        onSubmit={vi.fn()}
      />,
    );

    expect(screen.getByRole("button", { name: "Fold" })).toBeInTheDocument();
    expect(screen.queryByRole("button", { name: "Call 85" })).not.toBeInTheDocument();
    expect(screen.getByRole("button", { name: "All-in 76" })).toBeInTheDocument();
  });

  it("hides call when calling would be the same as going all-in", () => {
    render(
      <ActionBar
        roomId="room-001"
        pot={153}
        pendingAction={{
          token: "turn-1",
          seatIndex: 5,
          minAmount: 55,
          minBetAmount: 0,
          maxAmount: 55,
          canCheck: false,
          canCall: true,
          canBet: false,
          canFold: true,
          canAllIn: true,
        }}
        onSubmit={vi.fn()}
      />,
    );

    expect(screen.getByRole("button", { name: "Fold" })).toBeInTheDocument();
    expect(screen.queryByRole("button", { name: "Call 55" })).not.toBeInTheDocument();
    expect(screen.getByRole("button", { name: "All-in 55" })).toBeInTheDocument();
  });

  it("shows a live countdown when the pending action has an expiry", () => {
    vi.useFakeTimers();
    vi.setSystemTime(new Date("2026-04-25T00:00:00.000Z"));

    render(
      <ActionBar
        roomId="room-001"
        pot={8}
        pendingAction={{
          token: "turn-1",
          seatIndex: 5,
          minAmount: 2,
          minBetAmount: 4,
          maxAmount: 96,
          canCheck: false,
          canCall: true,
          canBet: true,
          canFold: true,
          canAllIn: true,
          expiresAt: Date.now() + 65000,
        }}
        onSubmit={vi.fn()}
      />,
    );

    expect(screen.getByText("65s")).toBeInTheDocument();
  });
});
