import type { RoomSnapshot } from "../lib/types";
import {
  ACTION_PLAYBACK_DELAY_MS,
  HAND_FINISH_PLAYBACK_DELAY_MS,
  ROUND_PLAYBACK_DELAY_MS,
  TURN_PLAYBACK_DELAY_MS,
  getPlaybackDelayMs,
} from "./roomPlayback";

function buildSnapshot(
  overrides: Partial<RoomSnapshot> = {},
): RoomSnapshot {
  return {
    roomId: "room-001",
    roomName: "Table 1",
    status: "running",
    viewerRole: "player",
    handNumber: 1,
    smallBlind: 1,
    pot: 3,
    currentAmount: 2,
    round: "PREFLOP",
    boardCards: [],
    seats: [],
    events: [],
    version: 1,
    ...overrides,
  };
}

describe("roomPlayback", () => {
  it("delays player actions, turns, rounds, and hand finish snapshots", () => {
    const base = buildSnapshot();

    expect(
      getPlaybackDelayMs(
        base,
        buildSnapshot({
          version: 2,
          events: [
            {
              kind: "player_action",
              message: "seat 1 called 2",
              handNumber: 1,
              round: "PREFLOP",
              seatIndex: 1,
              actionType: "CALL",
              amount: 2,
            },
          ],
        }),
      ),
    ).toBe(ACTION_PLAYBACK_DELAY_MS);

    expect(
      getPlaybackDelayMs(
        base,
        buildSnapshot({
          version: 2,
          status: "awaiting_action",
          events: [
            {
              kind: "turn",
              message: "seat 2 to act",
              handNumber: 1,
              round: "PREFLOP",
              seatIndex: 2,
            },
          ],
        }),
      ),
    ).toBe(TURN_PLAYBACK_DELAY_MS);

    expect(
      getPlaybackDelayMs(
        base,
        buildSnapshot({
          version: 2,
          round: "FLOP",
          events: [
            {
              kind: "round_start",
              message: "flop opened",
              handNumber: 1,
              round: "FLOP",
            },
          ],
        }),
      ),
    ).toBe(ROUND_PLAYBACK_DELAY_MS);

    expect(
      getPlaybackDelayMs(
        base,
        buildSnapshot({
          version: 2,
          status: "hand_finished",
          events: [
            {
              kind: "hand_finish",
              message: "hand 1 finished",
              handNumber: 1,
            },
          ],
        }),
      ),
    ).toBe(HAND_FINISH_PLAYBACK_DELAY_MS);
  });

  it("does not delay duplicate or non-animated snapshots", () => {
    const base = buildSnapshot();

    expect(getPlaybackDelayMs(base, base)).toBe(0);
    expect(
      getPlaybackDelayMs(
        base,
        buildSnapshot({
          version: 2,
          pot: 4,
        }),
      ),
    ).toBe(0);
  });
});
