

function grabPlayback(evt, elt) {
	const audioElt = htmx.find("#currently-playing")
	let grabbed = true
	elt.setPointerCapture(evt.pointerId)
	const parent = elt.parentNode
	let newProgress
	elt.onpointermove = (e) => {
		let relativeX = e.clientX - parent.offsetLeft
		relativeX = Math.max(0, relativeX)
		relativeX = Math.min(parent.offsetWidth, relativeX)
		newProgress = relativeX / parent.offsetWidth
		elt.style.left = `${relativeX - 5}px`
	}
	elt.onpointerup = (e) => {
		elt.releasePointerCapture(evt.pointerId)
		elt.onpointermove = null
		audioElt.currentTime = newProgress * audioElt.duration
		grabbed = false
	}
}

function getFullTime(totalSeconds) {
	const seconds = Math.floor(totalSeconds % 60)
	const minutes = Math.floor(totalSeconds / 60)
	return `${minutes.toString().padStart(2, '0')}:${seconds.toString().padStart(2, '0')}`
}




// function InitPlayer() {
// 	const init = () => {
// 		const totalMinutes = Math.floor(audio.duration / 60)
// 		const totalSeconds = Math.floor(audio.duration % 60)
// 		const totalMinutesElt = htmx.find('#player #duration #total-minutes')
// 		const totalSecondsElt = htmx.find('#player #duration #total-seconds')
// 		totalMinutesElt.innerHTML = totalMinutes.toString().padStart(2, '0')
// 		totalSecondsElt.innerHTML = totalSeconds.toString().padStart(2, '0')
// 		htmx.find("#page").style.marginBottom = '100px';
// 	}
// 	if (audio.readyState > 0) {
// 		init()
// 	} else {
// 		audio.onloadedmetadata = (e) => {
// 			init()
// 		}
// 	}
// }
//
// InitPlayer()


function HandleTimeClick(elt, e) {
	console.log('time click')
	const audioElt = htmx.find("#currently-playing")
	const parent = elt.parentNode
	let newProgress
	let relativeX = e.clientX - parent.offsetLeft
	relativeX = Math.max(0, relativeX)
	relativeX = Math.min(parent.offsetWidth, relativeX)
	newProgress = relativeX / parent.offsetWidth
	audioElt.currentTime = newProgress * audioElt.duration
}

function HandleMouseLeaveTime(elt, e) {
	const playbackTimeBar = htmx.find("#playback-time-bar")
	playbackTimeBar.style.height = '3px'
	playbackTimeBar.style.transform = 'translateY(-3.5px)'
	playbackTimeBar.style.borderTop = '0.5px solid lightgray'
	playbackTimeBar.style.borderBottom = ''
}

      // transform: translateY(-3.5px);
function HandleMouseoverTime(elt, e) {
	console.log("mouseover time")
	const playbackTimeBar = htmx.find("#playback-time-bar")
	const player = htmx.find("#player")
	playbackTimeBar.style.height = '6px'
	playbackTimeBar.style.transform = 'translateY(-6px)'
	playbackTimeBar.style.borderBottom = '0.5px solid lightgray'

	console.log("playbackTimeBar: ", playbackTimeBar)
	const audioElt = htmx.find("#currently-playing")
	const parent = elt.parentNode
	let newProgress
	let relativeX = e.clientX - parent.offsetLeft
	relativeX = Math.max(0, relativeX)
	relativeX = Math.min(parent.offsetWidth, relativeX)
	newProgress = relativeX / parent.offsetWidth
	const hovertime = Math.floor(newProgress * audioElt.duration)
	const newElt = document.createElement('div')
	newElt.id = "time-tooltip"
	newElt.style.left = `${e.clientX}px`
	newElt.innerText = getFullTime(hovertime)
	elt.appendChild(newElt)
}


function HandleMousemoveTime(elt, e) {
	const audioElt = htmx.find("#currently-playing")
	const parent = elt.parentNode
	const tooltipElt = document.getElementById('time-tooltip')
	let newProgress

	let relativeX = e.clientX - parent.offsetLeft
	relativeX = Math.max(0, relativeX)
	relativeX = Math.min(parent.offsetWidth, relativeX)
	newProgress = relativeX / parent.offsetWidth
	const newTooltipTime = Math.floor(newProgress * audioElt.duration)
	tooltipElt.innerText = getFullTime(newTooltipTime)
	tooltipElt.style.left = `${e.clientX}px`
}
