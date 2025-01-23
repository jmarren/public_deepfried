function downloadFromCard(elt) {
	console.log(elt)
	const cardElt = htmx.closest(elt, ".track-card-container")
	console.log(cardElt)
	const src = htmx.find(cardElt, ".play-button").dataset.audioSrc
	const download = document.createElement('a')
	download.setAttribute('download', src)
	download.click()
}

function MoveCarouselRight(elt) {
	const carouselElt = htmx.closest(elt, ".card-carousel")
	const carouselWidth = carouselElt.getBoundingClientRect().width 
	const cardsContainerElts = Array.from(htmx.findAll(carouselElt, ".cards-container"))
	let cardsContainerElt
	if (cardsContainerElts.length > 1) {
		cardsContainerElts[0].remove()
		cardsContainerElt = cardsContainerElts[1]
	} else {
		cardsContainerElt = cardsContainerElts[0]
	}
	const cardElt = htmx.find(cardsContainerElt, ".track-card-container")
	const cardWidth = cardElt.getBoundingClientRect().width
	const cardsArr = Array.from(htmx.findAll(cardsContainerElt, ".track-card-container"))
	const cardsContainerWidth = cardsContainerElt.getBoundingClientRect().width
	const slideDist = -1 * carouselWidth
	let currentTransform = parseInt(cardsContainerElt.dataset.currentTransform)
	const totalWidth = cardsArr.length * cardWidth
	const paddedWidth = totalWidth - carouselWidth

	if (Math.abs(currentTransform) > paddedWidth) { 
		const newCards = cardsContainerElt.cloneNode(true)
		cardsContainerElt.style.position = 'absolute';
		newCards.dataset.currentTransform = 0
		carouselElt.appendChild(newCards)
		htmx.process(newCards)
		
		const frames = [
			{ transform: `translateX(${carouselWidth}px)`},
			{ transform: `translateX(${20}px)`}
		];
		const animation = new KeyframeEffect(
			newCards,
			frames, 
			{ duration: 500, fill: 'forwards', easing: 'cubic-bezier(0.15, 0.83, 0.66, 1)'}
			)
		const player = new Animation(animation, document.timeline)
		player.play()
	} 

	const frames = [
		{ transform: `translateX(${currentTransform}px)`},
		{ transform: `translateX(${currentTransform + slideDist + cardWidth}px)`}
	];
	const animation = new KeyframeEffect(
		cardsContainerElt,
		frames, 
		{ duration: 500, fill: 'forwards', easing: 'cubic-bezier(0.15, 0.83, 0.66, 1)'}
		)
	const player = new Animation(animation, document.timeline)
	player.play()
	cardsContainerElt.dataset.currentTransform = currentTransform + slideDist
}


htmx.onLoad(function() {
	const allCards = htmx.findAll(".fade-in-slowly")
	for (let i = 0; i < allCards.length; i++) {
		setTimeout(() => {
			htmx.addClass(allCards[i], "fade-in-animation")
		}, 35 * i)
	}
})

function triggerCardAnimations(cardSection) {
	const cards = htmx.findAll(cardSection,".track-card-container")
	for (let i = 0;  i < cards.length; i++) {
	}
}

