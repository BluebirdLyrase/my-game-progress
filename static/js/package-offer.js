const prevButton = document.querySelector('.prev');
const nextButton = document.querySelector('.next');
const swiperWrapper = document.querySelector('.swiper-wrapper');
const cards = document.querySelectorAll('.swiper-card');
const cardQty = swiperWrapper.getAttribute('card-qty');
const cardWidth = cards[0].offsetWidth + 20;
const wrapperWidth = swiperWrapper.offsetWidth;
const maxCardOnScreen = Math.round(wrapperWidth / cardWidth) + 1;
// console.log('maxCardOnScreen:', maxCardOnScreen);
let currentIndex = 0;

function updateSwiperPosition() {
  swiperWrapper.style.transform = `translateX(-${currentIndex * cardWidth}px)`;
}

prevButton.addEventListener('click', () => {
  if (currentIndex > 0) {
    currentIndex--;
  }
  updateSwiperPosition();
});

nextButton.addEventListener('click', () => {
  if (currentIndex < cards.length - maxCardOnScreen) {
    currentIndex++;
    updateSwiperPosition();
  }
});

function onClickBuy(param) {
  //TODO
}
