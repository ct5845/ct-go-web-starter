Alpine.store("bottomtabs", {
  visible: true,
  hide() {
    this.visible = false;
  },
  show() {
    this.visible = true;
  },
});

Alpine.data("bottomtabs", () => {
  return {
    get visible() {
      return Alpine.store("bottomtabs").visible;
    },
  };
});
