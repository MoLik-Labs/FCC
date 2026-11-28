(() => {
  const readingTheme = document.getElementById("reading-theme");
  if (readingTheme) {
    const allowed = new Set(["light", "forest", "catppuccin-mocha"]);
    let saved = "light";
    try {
      saved = localStorage.getItem("fcc-docs-theme") || "light";
    } catch {
      /* Storage is optional. */
    }
    const apply = (value) => {
      const theme = allowed.has(value) ? value : "light";
      document.documentElement.dataset.readingTheme = theme;
      readingTheme.value = theme;
    };
    apply(saved);
    readingTheme.parentElement.hidden = false;
    readingTheme.addEventListener("change", () => {
      apply(readingTheme.value);
      try {
        localStorage.setItem("fcc-docs-theme", readingTheme.value);
      } catch {
        /* Keep the current page usable. */
      }
    });
  }
  const toggle = document.querySelector(".menu-toggle");
  const navigation = document.getElementById("navigation");
  if (!toggle || !navigation) return;
  const setOpen = (open) => {
    navigation.classList.toggle("open", open);
    toggle.setAttribute("aria-expanded", String(open));
  };
  toggle.addEventListener("click", () =>
    setOpen(toggle.getAttribute("aria-expanded") !== "true"),
  );
  navigation.addEventListener("click", (event) => {
    if (event.target.closest("a")) setOpen(false);
  });
  document.addEventListener("keydown", (event) => {
    if (
      event.key === "Escape" &&
      toggle.getAttribute("aria-expanded") === "true"
    ) {
      setOpen(false);
      toggle.focus();
    }
  });
})();
