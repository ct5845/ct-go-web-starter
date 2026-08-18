Alpine.data("combobox", (id, multi) => ({
  selectedText: "",
  loadingNextPage: false,
  resultsId: id + "-results",

  init() {
    this.selectedText = this.$el.dataset.selectedText;

    // subtree:true so this keeps working across loadNextPage's
    // outerHTML swap, which replaces the #resultsId-page node itself
    // rather than mutating it in place.
    const results = document.getElementById(this.resultsId);
    new MutationObserver(() => this.dropLoadedPins()).observe(results, {
      childList: true,
      subtree: true,
    });
  },

  loadNextPage(event) {
    if (this.loadingNextPage) return;

    const results = event.target;
    const nearBottom = results.scrollTop + results.clientHeight >= results.scrollHeight - 48;
    if (!nearBottom) return;

    const page = results.querySelector("[data-next-page-href]");
    const href = page.dataset.nextPageHref;
    if (!href) return;

    this.loadingNextPage = true;
    htmx.ajax("GET", href, {
      target: "#" + page.id,
      select: "#" + page.id,
      swap: "outerHTML",
    }).then(() => {
      this.loadingNextPage = false;
    });
  },

  // A pinned row (a selected item not yet on any loaded page) duplicates
  // the real row once its page loads — same name/value, so checking one
  // natively unchecks the other, but both would stay visible. Drop the
  // pinned copy once the real row is on the page.
  dropLoadedPins() {
    const pinned = document.getElementById(this.resultsId + "-pinned");
    if (!pinned) return;

    const loadedValues = new Set(
      Array.from(document.querySelectorAll("#" + this.resultsId + "-page input")).map((input) => input.value),
    );
    pinned.querySelectorAll(".combobox-option").forEach((option) => {
      const input = option.querySelector("input");
      if (loadedValues.has(input.value)) {
        option.remove();
      }
    });
  },

  focusSearch(event) {
    if (event.newState !== "open") return;
    this.$refs.search.focus();
    this.$refs.search.select();
  },

  pick(event) {
    if (!multi) {
      document.getElementById(id + "-popover").hidePopover();
    }
    this.refreshSelectedText();
    this.$dispatch("field-change", { id, value: event.target.value });
  },

  refreshSelectedText() {
    const popover = document.getElementById(id + "-popover");
    const checked = popover.querySelectorAll("input:checked");
    this.selectedText = Array.from(checked)
      .map((input) => input.dataset.label)
      .join(", ");
  },
}));
