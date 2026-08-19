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
    new MutationObserver(() => this.dropLoadedHiddenOptions()).observe(results, {
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

  // A hidden option (a selected item not yet on any loaded page) submits
  // the same name/value as its eventual visible row. Once that row loads,
  // drop the hidden one so a later uncheck of the visible row can't be
  // silently overridden by the hidden duplicate still being checked.
  dropLoadedHiddenOptions() {
    const loadedValues = new Set(
      Array.from(document.querySelectorAll("#" + this.resultsId + "-page input")).map((input) => input.value),
    );
    document.querySelectorAll("#" + this.resultsId + " [data-hidden-option]").forEach((input) => {
      if (loadedValues.has(input.value)) {
        input.remove();
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
