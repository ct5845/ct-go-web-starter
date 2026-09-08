const top = document.getElementById("<<< .ID >>>-search-top");
const bottom = document.getElementById("<<< .ID >>>-search-bottom");

if (top && bottom) {
  top.addEventListener("input", () => {
    bottom.value = top.value;
  });
  bottom.addEventListener("input", () => {
    top.value = bottom.value;
  });
}
