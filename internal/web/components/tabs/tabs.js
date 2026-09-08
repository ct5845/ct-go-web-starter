const strip = document.currentScript.previousElementSibling;
const active = strip.querySelector('[aria-current="page"]');
if (active) {
  active.scrollIntoView({ block: "nearest", inline: "center" });
}
