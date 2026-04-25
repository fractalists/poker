import type { RoomEvent, RoomSnapshot } from "../lib/types";

export const ROUND_PLAYBACK_DELAY_MS = 260;
export const TURN_PLAYBACK_DELAY_MS = 360;
export const ACTION_PLAYBACK_DELAY_MS = 480;
export const HAND_FINISH_PLAYBACK_DELAY_MS = 1080;

function getNewEvents(
  previous: RoomSnapshot | null,
  next: RoomSnapshot,
): RoomEvent[] {
  const previousEvents = previous?.events ?? [];
  const nextEvents = next.events ?? [];
  if (nextEvents.length <= previousEvents.length) {
    return [];
  }
  return nextEvents.slice(previousEvents.length);
}

function getLatestAnimatedEvent(
  previous: RoomSnapshot | null,
  next: RoomSnapshot,
): RoomEvent | undefined {
  const newEvents = getNewEvents(previous, next);
  return newEvents.at(-1);
}

export function getPlaybackDelayMs(
  previous: RoomSnapshot | null,
  next: RoomSnapshot,
): number {
  const previousVersion = previous?.version ?? -1;
  const nextVersion = next.version ?? previousVersion;

  if (nextVersion <= previousVersion) {
    return 0;
  }

  const latestEvent = getLatestAnimatedEvent(previous, next);
  if (!latestEvent) {
    return 0;
  }

  switch (latestEvent.kind) {
    case "player_action":
      return ACTION_PLAYBACK_DELAY_MS;
    case "turn":
      return TURN_PLAYBACK_DELAY_MS;
    case "round_start":
      return ROUND_PLAYBACK_DELAY_MS;
    case "hand_finish":
      return HAND_FINISH_PLAYBACK_DELAY_MS;
    default:
      return 0;
  }
}
