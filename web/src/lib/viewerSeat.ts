export function viewerSeatKey(roomId: string): string {
  return `poker.viewerSeat.${roomId}`;
}

export function viewerTokenKey(roomId: string): string {
  return `poker.viewerToken.${roomId}`;
}
