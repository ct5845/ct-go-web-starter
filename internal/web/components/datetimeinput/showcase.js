Alpine.data("showcaseDateRange", (startId, endId) => ({
  previousStart: null,

  init() {
    this.previousStart = document.getElementById(startId).value;
  },

  shiftEnd(event) {
    if (event.detail.id !== startId) return;

    const deltaDays = (new Date(event.detail.value) - new Date(this.previousStart)) / 86400000;
    this.previousStart = event.detail.value;
    if (!Number.isFinite(deltaDays) || deltaDays === 0) return;

    const end = document.getElementById(endId);
    const shifted = new Date(end.value);
    shifted.setDate(shifted.getDate() + deltaDays);
    end.value = shifted.toISOString().slice(0, 10);
  },
}));
