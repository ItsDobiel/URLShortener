document.addEventListener("DOMContentLoaded", function () {
  var currentSlide = 0;
  var slides = document.querySelectorAll(".slide");
  var keyboardHintTimeout;
  slides.forEach(function (slide, index) {
    if (index > 0) {
      var slideNum = index + " / " + (slides.length - 1);
      slide.setAttribute("data-slide-number", slideNum);
    }
  });

  function showSlide(n) {
    if (slides.length === 0) return;
    slides[currentSlide].classList.remove("active");
    if (n >= slides.length) {
      currentSlide = slides.length - 1;
    } else if (n < 0) {
      currentSlide = 0;
    } else {
      currentSlide = n;
    }
    slides[currentSlide].classList.add("active");
    slides[currentSlide].scrollTop = 0;
    updateNavigation();
    updateProgressBar();
    showKeyboardHint();
  }

  window.changeSlide = function (direction) {
    showSlide(currentSlide + direction);
  };

  function updateNavigation() {
    var prevBtn = document.getElementById("prevBtn");
    var nextBtn = document.getElementById("nextBtn");
    if (prevBtn) prevBtn.disabled = currentSlide === 0;
    if (nextBtn) nextBtn.disabled = currentSlide === slides.length - 1;
  }

  function updateProgressBar() {
    var bar = document.getElementById("progressBar");
    if (!bar) return;
    var total = slides.length > 1 ? slides.length - 1 : 1;
    var progress = (currentSlide / total) * 100;
    bar.style.width = progress + "%";
  }

  function showKeyboardHint() {
    var hint = document.getElementById("keyboardHint");
    if (!hint) return;
    hint.classList.add("show");
    clearTimeout(keyboardHintTimeout);
    keyboardHintTimeout = setTimeout(function () {
      hint.classList.remove("show");
    }, 3000);
  }

  document.addEventListener("keydown", function (e) {
    switch (e.key) {
      case "ArrowLeft":
      case "ArrowUp":
      case "PageUp":
        changeSlide(-1);
        break;
      case "ArrowRight":
      case "ArrowDown":
      case "PageDown":
      case " ":
        e.preventDefault();
        changeSlide(1);
        break;
      case "Home":
        showSlide(0);
        break;
      case "End":
        showSlide(slides.length - 1);
        break;
      case "f":
      case "F":
        if (document.fullscreenElement) {
          document.exitFullscreen();
        } else {
          document.documentElement.requestFullscreen();
        }
        break;
      case "?":
        showKeyboardHint();
        break;
    }
  });

  var touchStartX = 0;
  var touchEndX = 0;

  document.addEventListener("touchstart", function (e) {
    touchStartX = e.changedTouches[0].screenX;
  });

  document.addEventListener("touchend", function (e) {
    touchEndX = e.changedTouches[0].screenX;
    handleSwipe();
  });

  function handleSwipe() {
    if (touchEndX < touchStartX - 50) changeSlide(1);
    if (touchEndX > touchStartX + 50) changeSlide(-1);
  }

  updateNavigation();
  updateProgressBar();

  setTimeout(function () {
    showKeyboardHint();
  }, 500);
});
