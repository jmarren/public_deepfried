

function toggleBpmType(elt) {
	const useExactCbx = htmx.find("input[value=use-exact]")
	const useRangeCbx = htmx.find("input[value=use-range]")
	const exactInput = htmx.find("input[name=exact-bpm]")
	const minBpmInput = htmx.find("input[name=min-bpm]")
	const maxBpmInput = htmx.find("input[name=max-bpm]")
	const minBpmLabel = htmx.find("label[for=min-bpm]")
	const maxBpmLabel = htmx.find("label[for=max-bpm]")
	const minMaxContainer = htmx.find("#bpm-min-max-inputs-container")
	const displaySlider = htmx.find("#display-slider")
	
	if (useExactCbx.checked == true) {
		hideRangeElts()
		showExactElts()
		minBpmInput.disabled = true;
		maxBpmInput.disabled = true;
		minBpmLabel.style.color = ''
		maxBpmLabel.style.color = ''
		exactInput.disabled = false;
		useRangeCbx.checked = false;
	} else if (useRangeCbx.checked == true) {
		hideExactElts()
		showRangeElts()
		minBpmInput.disabled = false;
		maxBpmInput.disabled = false;
		minBpmLabel.style.color = 'var(--dark-gray)'
		maxBpmLabel.style.color = 'var(--dark-gray)'
		exactInput.disabled = true;
	}

	if (!useRangeCbx.checked) {
		hideRangeElts()
		minBpmInput.disabled = true;
		maxBpmInput.disabled = true;
		minBpmLabel.style.color = ''
		maxBpmLabel.style.color = ''
	} 
	if (!useExactCbx.checked) {
		hideExactElts()
		exactInput.disabled = true;
	}
}


function hideExactElts() {
	console.log('hiding exact elts')
	const elts = htmx.findAll(".bpm-exact") 
	elts.forEach((elt) => {
		console.log(elt)
		elt.style.visibility = 'hidden'
	})
}

function showExactElts() {
	const elts = htmx.findAll(".bpm-exact") 
	elts.forEach((elt) => {
		console.log(elt)
		elt.style.visibility = 'visible'
	})
}
function showRangeElts() {
	const rangeElts = htmx.findAll(".bpm-range") 
	rangeElts.forEach((elt) => {
		elt.style.visibility = 'visible'
	})
}


function hideRangeElts() {
	console.log('hiding range elts')
	const rangeElts = htmx.findAll(".bpm-range") 
	console.log("rangeElts: ", rangeElts)
	rangeElts.forEach((elt) => {
		elt.style.visibility = 'hidden'
	})
}



function grabSlider(elt, evt) {
	console.log(elt,evt)
	const useRangeCbx = htmx.find("input[value=use-range]")
	if (!useRangeCbx.checked) return
	htmx.find("input[name=exact-bpm]").disabled = true;
	const MAX_ALLOWED_BPM = 200;
	const SLIDER_WIDTH = htmx.find("#display-slider").getBoundingClientRect().width;
		
	
	let otherSlider
	let inputElt
	if (elt.id == "low-slider") { 
		inputElt = htmx.find("input[name=min-bpm]")
		otherSlider = htmx.find("#high-slider")
	} else if (elt.id == "high-slider"){
		inputElt = htmx.find("input[name=max-bpm]")
		otherSlider = htmx.find("#low-slider")
	}

	let grabbed = true
	let pointerId = elt.setPointerCapture(evt.pointerId)
	const parent = htmx.find("#display-slider")
	const parentRect = parent.getBoundingClientRect()
	const displayLeft = parentRect.left
	const displayRight = parentRect.right
	const displayWidth = parentRect.width
	const eltWidth = elt.getBoundingClientRect().width
	const eltRad = eltWidth / 2
	let newProgress
	const maxLeft = SLIDER_WIDTH - eltRad
	const otherSliderLeft = otherSlider.getBoundingClientRect().left - displayLeft
	console.log("otherSliderLeft: ", otherSliderLeft)

	window.onpointerup = (evt) => {
		if (!grabbed) return
		elt.releasePointerCapture(pointerId)
		grabbed = false
		window.onpointerup = null;
	 }
	

	window.onpointermove = (evt) => {
		if (!grabbed) return
		console.log(evt)

		// get the clientX of event and subtract the offset from window left to the #display-slider, then subtract the radius of the slider
		let  newLeft = evt.clientX - displayLeft - eltRad
		console.log(newLeft)
		
		if (evt.clientX > displayRight) {
			newLeft = displayWidth - eltRad
		} else if (evt.clientX < (displayLeft - eltRad)) {
			newLeft =  -1 * eltRad
		} 
	
		if (elt.id == "low-slider") {
			if (newLeft >= (otherSliderLeft)) {
				newLeft = otherSliderLeft
			} 
		} else if (elt.id == "high-slider") {
			if (newLeft <= (otherSliderLeft)) {
				newLeft = otherSliderLeft
			} 
		}

		let newInputVal = ( (newLeft + eltRad) / displayWidth) * MAX_ALLOWED_BPM
		if (newInputVal < 0) {
			newInputVal = 0
		}
		inputElt.value = Math.round(newInputVal)
		elt.style.left = `${newLeft}px`
			
		let lowSliderLeft
		let highSliderLeft
		if (elt.id == "low-slider") {
			lowSliderLeft = newLeft + eltRad
			highSliderLeft = otherSliderLeft + eltRad
		} else if (elt.id == "high-slider") {
			lowSliderLeft = otherSliderLeft + eltRad
			highSliderLeft = newLeft + eltRad
		}
		
		parent.style.background = `linear-gradient(90deg,black ${lowSliderLeft}px, white ${lowSliderLeft}px, white ${highSliderLeft}px, black ${highSliderLeft}px)`
	} 
}
		/*
		// let relativeX = evt.clientX - parent.offsetLeft
		// relativeX = Math.max(0, relativeX)
		// relativeX = Math.min(parent.offsetWidth, relativeX)
		// newProgress = relativeX / parent.offsetWidth
		// elt.style.left = `${relativeX - 5}px`
		//
		// console.log("relativeX: ", relativeX);
		// const relativePercent = SLIDER_WIDTH / relativeX
		//
		// console.log(relativePercent)

		if (otherSliderLeft > relativeX) {
			document.documentElement.style.setProperty('--bpm-slider-min',`${relativeX - 5}px`)
			document.documentElement.style.setProperty('--bpm-slider-max',`${otherSliderLeft}px`)
			htmx.find("input[name=min-bpm]").value = relativeX  * COEFFICIENT
			htmx.find("input[name=max-bpm]").value = otherSliderLeft * COEFFICIENT
		} else {
			document.documentElement.style.setProperty('--bpm-slider-max',`${relativeX - 5}px`)
			document.documentElement.style.setProperty('--bpm-slider-min',`${otherSliderLeft}px`)
			htmx.find("input[name=min-bpm]").value = otherSliderLeft * COEFFICIENT
			htmx.find("input[name=max-bpm]").value = relativeX * COEFFICIENT
		}
	}
	window.onpointerup = () => {
		elt.releasePointerCapture(e.pointerId)
		elt.onpointermove = null
		grabbed = false
		window.onpointerup = null;
	}
	*/


/*
function grabBpmSlider(elt, e) {
	const useRangeCbx = htmx.find("input[value=use-range]")
	if (!useRangeCbx.checked) return
	htmx.find("input[name=exact-bpm]").disabled = true;
	const MAX_ALLOWED_BPM = 200;
	const SLIDER_WIDTH = htmx.find("#bpm-slider-bar").getBoundingClientRect().width;
	const COEFFICIENT =  MAX_ALLOWED_BPM / SLIDER_WIDTH;

	let otherSlider;
	if (elt.id == "bpm-slider-one") {
		otherSlider = htmx.find("#bpm-slider-two")
	} else {
		otherSlider = htmx.find("#bpm-slider-one")
	}

	
	let otherSliderLeft  = parseInt(otherSlider.style.left.trimEnd().trimEnd());
	if (!otherSliderLeft) {
		otherSliderLeft = 0;
	}


	let grabbed = true
	elt.setPointerCapture(e.pointerId)
	const parent = elt.parentNode
	let newProgress
	elt.onpointermove = (evt) => {
		let relativeX = evt.clientX - parent.offsetLeft
		relativeX = Math.max(0, relativeX)
		relativeX = Math.min(parent.offsetWidth, relativeX)
		newProgress = relativeX / parent.offsetWidth
		elt.style.left = `${relativeX - 5}px`
		
		
		
		console.log("relativeX: ", relativeX, "\notherSliderLeft: ", otherSliderLeft);

		if (otherSliderLeft > relativeX) {
			document.documentElement.style.setProperty('--bpm-slider-min',`${relativeX - 5}px`)
			document.documentElement.style.setProperty('--bpm-slider-max',`${otherSliderLeft}px`)
			htmx.find("input[name=min-bpm]").value = relativeX  * COEFFICIENT
			htmx.find("input[name=max-bpm]").value = otherSliderLeft * COEFFICIENT
		} else {
			document.documentElement.style.setProperty('--bpm-slider-max',`${relativeX - 5}px`)
			document.documentElement.style.setProperty('--bpm-slider-min',`${otherSliderLeft}px`)
			htmx.find("input[name=min-bpm]").value = otherSliderLeft * COEFFICIENT
			htmx.find("input[name=max-bpm]").value = relativeX * COEFFICIENT
		}
	}
	window.onpointerup = () => {
		elt.releasePointerCapture(e.pointerId)
		elt.onpointermove = null
		grabbed = false
		window.onpointerup = null;
	}

}
*/
