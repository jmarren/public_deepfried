
function RemoveModal(event) {
    console.log("remove modal event: ", event)
    const modal = htmx.find('#modal')
    if (modal.contains(event.explicitOriginalTarget)) return
    modal.innerHTML = ''
    const currentUrl = window.location.href
    const newUrl = new URL(currentUrl)
    newUrl.searchParams.delete("modal")
    window.history.replaceState({}, null, newUrl)
}



window.workers = []

const PAGES = ["explore", "user-feed", "search-results-page", "my-profile", "my-downloads", "my-uploads"]
const SIDEBAR_PAGES = ["explore", "user-feed", "submissions"]
const SIDEBAR_BUTTONS = ["#globe-button", "#home-button", "#submissions-button"]

////////////////////////////////////////////////////////////// hx-audio extension
var hxAPI 

  // When the audio time is updated, update the player. 
  // if the audio is almost over, trigger a preload event for the next track
  // by triggering a mouseover on the skip forward button
  function audioTimeupdateListener(event) {
      const audioElt = event.target
      updatePlayer(audioElt)
      const currentCanvas =  hxAPI.getCurrentCanvas()
      if (currentCanvas) {
	  currentCanvas.timeUpdate(audioElt.currentTime / audioElt.duration)
      }

      // const progress = event.target.currentTime / evet
      if ((audioElt.duration - audioElt.currentTime) < 4) { 
	  const skipforwardButton = document.querySelector("#skip-forward")
	  const event = new MouseEvent('mouseover', {
	       bubbles: true,
	       cancelable: true,
	       view: window
	     });
	 skipforwardButton.dispatchEvent(event)
      }
  }

function updatePlayer(elt) {
        const subElts = hxAPI.playerSubElts
	if (!subElts) return
	const playerMinutes = subElts.minutes
	const playerSeconds = subElts.seconds
	const progressMarker = subElts.progressMarker
	const playbackTimeBar = subElts.playbackTimeBar

	if (!playerMinutes || !playerSeconds || !progressMarker) {
	  return;
	}
	const seconds = Math.floor(elt.currentTime % 60)
	const minutes = Math.floor(elt.currentTime / 60)
	playerSeconds.innerHTML = seconds.toString().padStart(2, '0')
	playerMinutes.innerHTML = minutes.toString().padStart(2, '0')
	const progress = elt.currentTime / elt.duration
	progressMarker.style.left = `calc(${progress * 100}% - 5px)`
	playbackTimeBar.style.background = `linear-gradient(90deg,black 0%, black ${progress * 100}%, white ${progress *  100}%, white ${100}%)`
}

function HandleAudioMetadataLoaded(audioElt) {
	console.log("metadata loaded")
        subElts = hxAPI.playerSubElts
	if (!subElts) {
	  console.log('subElts not defined')
	  return
	}
	const totalMinutes = subElts.totalMinutes
        const totalSeconds = subElts.totalSeconds
	if (!totalMinutes || !totalSeconds) {
	    console.log('subElts not defined')
	    return
	  }

	const seconds = Math.floor(audioElt.duration % 60)
	const minutes = Math.floor(audioElt.duration / 60)
	totalSeconds.innerHTML = seconds.toString().padStart(2, '0')
	totalMinutes.innerHTML = minutes.toString().padStart(2, '0')
  }



htmx.defineExtension('audio', {
 onEvent: function(name, event) {

    /* ------------- oobBeforeSwap  ------------- */

      if (name === "htmx:oobBeforeSwap") {
	// const element = event.target || event.detail.elt;

      }


      /*--------------- oobAfterSwap -----------*/
      if (name === "htmx:oobAfterSwap") {
	const element = event.target || event.detail.elt;

	// if the new node is a #currently-playing audio elt, add an event listener to 
        // preload the next short queue when it is almost over
	// if the processed node is the <audio> element 
	// set it with hxAPI only if:
	//    1. hxAPI.audioElt is not yet defined 
	// or 2. the src of the new element does not match hxAPI audioElt's current 
	//	 src

	if (element.id === "currently-playing") {
	  if (!hxAPI.audioElt) {
	      hxAPI.setAudioElt(element)
	    }  else if (element.src != hxAPI.audioElt.src) {
	      hxAPI.setAudioElt(element)
	  }
	}
      }

      /*--------------- beforeCleanupElement -----------*/
     if (name == 'htmx:beforeCleanupElement') {
	const parent = event.target || event.detail.elt;

	const hxAudioElts = htmx.findAll(parent, "[hx-audio]")
	if (hxAudioElts) {
	    hxAPI.removeHxAudioElts(Array.from(hxAudioElts))
	}

	const toggleButtons = htmx.findAll(parent, "[hx-toggle-audio]")
	if (toggleButtons) {
	  hxAPI.removeToggleButtons(toggleButtons)
	}
     }

      /*--------------- beforeSwap -----------*/

     if (name === "htmx:beforeSwap") {
    

	    // if the element to swap is the player, morph it into the existing player with idiomorph
	    // then process the new node with htmx.process after it is added
	    if (event.target.id === "player") {
		if (hxAPI.playerElt) {
		Idiomorph.morph(event.detail.elt, event.detail.serverResponse)
		  }
	      }
	  }


      /*--------------- afterProcessNode -----------*/
     if (name === 'htmx:afterProcessNode') {
      const parent = event.target || event.detail.elt;
	

	const modal = htmx.find("#modal")




	if (PAGES.includes(parent.id)) {
		if (hxAPI.playerElt) {
		    htmx.find("#page-content").style.paddingBottom = '100px'    
		}
		const pageIndex = SIDEBAR_PAGES.indexOf(parent.id) 
	    
		SIDEBAR_PAGES.forEach((page, index) => {
			const currentButton =  htmx.find(SIDEBAR_BUTTONS[index])
			if (index == pageIndex) {
			    htmx.addClass(currentButton, "active-page-button")
			    htmx.removeClass(currentButton, "inactive-page-button")
			} else {
			    htmx.removeClass(currentButton, "active-page-button")
			    htmx.addClass(currentButton, "inactive-page-button")
			}
		})
	} 


 
	// if the processed node is the player (and it has innerHTML)
	// set it with hxAPI only if:
	//    1. hxAPI.playerElt is not defined yet
	// or 2. the innerHTML of the new element does not match hxAPI's current 
	//	 innerHTML
	const player = htmx.find("#player") 
	if (player && player.innerHTML ) {
           if (!hxAPI.playerElt) {
		hxAPI.setPlayerElt(player)
	    } else if (player.innerHTML != hxAPI.playerElt.innerHTML) {
		hxAPI.setPlayerElt(player)
	    }
	}


      
	// if the processed element is the contains hx-audio-vis canvases, add them 
	const canvases = htmx.findAll(parent, '[hx-audio-vis]')
	if (canvases.length > 0) {
	  hxAPI.addCanvases(canvases)
	}
      

	// add any [hx-audio] elements to the api
	const hxAudioElts = htmx.findAll(parent, "[hx-audio]")
	if (hxAudioElts.length > 0) {
	    const eltsToAdd = Array.from(hxAudioElts)
	    hxAPI.addHxAudioElts(eltsToAdd)
	 }
      }

      /*--------------- configRequest -----------*/
       if (name === "htmx:configRequest") {
        const elt = event.target || event.detail.elt;

	// if the request is for the player, get the current from the queue and append the 
	//  appropriate form data to the request
	if (event.detail.path === "/player") {
	      let playingId = hxAPI.getCurrentFromQueue() 
	      const queue = hxAPI.getNextThree()
	      event.detail.formData.append("playing", playingId )
	      event.detail.formData.append("queue", queue )
	      if (hxAPI.playerElt) {
		    event.detail.headers["X-Player-Present"]  = "true"
	     }
        }

	const reqPath = event.detail.path
	if (reqPath.includes("modal") && event.target.id != "file-input" ) {
		const modal = htmx.find("#modal")
		if (modal.innerHTML) {
		    event.preventDefault()
		    modal.innerHTML =''
		}
	}

        // always add the 'HX-Audio-Queue' and 'HX-Audio-Playing' headers to htmx requests
	  event.detail.headers['HX-Audio-Queue'] = hxAPI.getNextThree()
	  event.detail.headers['HX-Audio-Playing'] = hxAPI.getPlayingId()
	}




        return
  },
  init: function(api) {
    
    api.audioQueue = []
    api.audioPlaying = ""
    api.globalQueue = []
    api.globalQueueIndex = 0
    api.templateMap = new Map()
    api.templateMap.set(
	"pausebars",  
        `<div class="pause-bars">
            <div class="pause-bar"></div>
            <div class="pause-bar"></div>
          </div>`
    )
    api.templateMap.set(
	"playtriangle",
	`<div class="play-triangle"></div>`
    )
    api.audioElt = null
    api.toggleButtons = new Set()
    api.hxAudioElts = []
    api.playerElt = null
    api.playerSubElts = {}
    api.canvases = new Set()
    api.currentCanvase = null
    api.pageElt = htmx.find("#page-content")
    



    // toggle buttons ///////////////////////
    api.addToggleButtons = (elts) => {
      elts.forEach((elt) => {
	api.toggleButtons.add(elt)
        htmx.on(elt, 'click', api.toggle)
      })
    }

    api.removeToggleButtons = (elts) => {
	elts.forEach((elt) => {
	    api.toggleButtons.delete(elt)
	})
    }

    api.getToggleButtons = () => {
	return api.toggleButtons
    }

    // hx-audio elts ///////////////////////////
    api.addHxAudioElts = (elts) => {
	elts.forEach((elt) => {
	  htmx.on(elt, 'click', function(event) {
	    api.playFromElt(elt, event)
	  })
	  if (!api.hxAudioElts.includes(elt)) {
	      api.hxAudioElts.push(elt)
	  }
	})
    }

    api.removeHxAudioElts = (elts) => {
	api.hxAudioElts = api.hxAudioElts.filter((elt) => {
	    return !elts.includes(elt)
	})
    }
    
    api.getHxAudioElts = () => {
	return api.hxAudioElts
    }

    api.getHxAudioEltsWithId = (audioId) => {
	return api.hxAudioElts.filter((elt) => {
	    return elt.getAttribute('hx-audio') === audioId
	})
    }

    // <audio> elt /////////////////////////////////////////
    api.setAudioElt = (audioElt) => {
	api.audioElt = audioElt

	htmx.on(audioElt, 'htmx:afterSettle', function(e) {
	    api.play()
	})
	api.updateCurrentCanvas()
	htmx.on(audioElt, 'timeupdate', audioTimeupdateListener)
    }

    api.getAudioElt = () => {
	return api.audioElt
    }

    // <audio> elt state ///////////////////////
    api.isPaused = () => {
      const audioElt = api.getAudioElt()
      if (!audioElt) return
      return audioElt.paused
    }

    api.setAudioEltAudioId = (audioId) => {
      const audioElt = api.getAudioElt()
      if (!audioElt) return
      audioElt.dataset.audioId = audioId
    }

    api.pause = () =>  {
      const audioElt = api.getAudioElt() 
      if (!audioElt) return
      audioElt.pause()
      api.toggleAllButtons()
    }

    api.play = () => {
      const audioElts = api.getHxAudioElts()
      const audioElt = api.getAudioElt() 
      if (!audioElt) return
      audioElt.play()
      api.toggleAllButtons()
    }
  
    api.toggle = () => {
	if (api.isPaused()) {
	  api.play()
	} else {
	  api.pause()
	}
    }

    // player elt ///////////////////////////////////////
    api.setPlayerElt = (player) => {
	    if (!api.playerElt) {
		    htmx.addClass(player, "slide-up")
	    }
	
	    api.playerElt = player
	    htmx.process(player)


	    const toggleButtons = htmx.findAll(player, "[hx-toggle-audio]")
	    if (toggleButtons.length > 0) {
		api.addToggleButtons(Array.from(toggleButtons))
	    }
	    const audioForward = htmx.find(player, "[hx-audio-forward]")
	    if (audioForward) {
		audioForward.addEventListener('click', function() {
		    api.skipForward()
		})
	    }
	    const audioBackward = htmx.find(player, "[hx-audio-backward]")
	    if (audioBackward) {
		audioBackward.addEventListener('click', function() {
		    api.skipBackward()
		})
	    }

	
	api.playerSubElts = {
	  minutes: htmx.find(player, "#current-time  #current-minutes"),
	  seconds: htmx.find(player, "#current-time #current-seconds"),
	  progressMarker: htmx.find(player, "#playback-time-bar #playback-progress-marker"),
	  playbackTimeBar: htmx.find(player, "#playback-time-bar"),
	  totalMinutes: htmx.find(player, "#duration #total-minutes"),
	  totalSeconds: htmx.find(player, "#duration #total-seconds"),
	}

	if (api.audioElt.readyState > 0) {
		    console.log('readyState > 0')
	            HandleAudioMetadataLoaded(api.audioElt)
	      } else {
		console.log('adding loadedmetadata event listener')
		htmx.on(api.audioElt, 'loadedmetadata', function(e) {
			HandleAudioMetadataLoaded(api.audioElt)
		})
	}

	htmx.find("#page-content").style.paddingBottom = '100px'
    }

    api.getAllCurrentIdButtons = () => {
	return api.getHxAudioEltsWithId(api.getCurrentId())
    }

    api.toggleAllButtons = () => {
      const toggleButtons = api.getToggleButtons()
      const hxAudioElts = api.getHxAudioElts()
      const playingId = api.getPlayingId()

      // TODO could store this as currentHxAudioElts 
      hxAudioElts.forEach((button) => {
	  const buttonId = button.getAttribute("hx-audio")
	  if (button.hasAttribute('hx-audio-toggle-fx')) { 
	      const isPlaying = (buttonId === playingId)
	      api.toggleAudioFx(button, isPlaying)
	  }
	})

       toggleButtons.forEach((button) =>{
	      api.toggleAudioFx(button, !api.isPaused())
       })
    }

    api.toggleAudioFx = (elt, isPlaying) => {
	toggleFx = elt.getAttribute('hx-audio-toggle-fx').trim()
	if (!toggleFx) return
	toggleFx = toggleFx.trim()
        const expressionArr = toggleFx.match(/\S+/g) || []

	let toggleTarget = elt.getAttribute("hx-audio-toggle-target")
	
	// ie hx-audio-toggle-fx="class id-is-playing"
	if (expressionArr[0] == "class") {
	    let targetElt = elt
	    if (toggleTarget) {
		const targetExpressionArr = toggleTarget.match(/\S+/g) || []
		if (targetExpressionArr[0] === "closest") {
		  toggleTarget = htmx.closest(elt, targetExpressionArr[1])
		} else {
		  toggleTarget = targetExpressionArr[0]
		}
	    }
	  if (isPlaying) {
	    htmx.addClass(toggleTarget, expressionArr[1])
	  } else {
	    htmx.removeClass(toggleTarget, expressionArr[1])
	  }
	    return
	}
        
	// ie hx-audio-toggle-fx="innerHTML playtriangle pausebars"
	if (["innerHTML" || "outerHTML" || "beforebegin" || "afterend"].includes(expressionArr[0])) {  
	    let swapContent =""
	     if (isPlaying) {
		swapContent = api.templateMap.get(expressionArr[2])
	     } else {
		swapContent = api.templateMap.get(expressionArr[1])
	      }
	      switch (expressionArr[0]) {
		case 'innerHTML': 
		    elt.innerHTML = swapContent;
		  break;
		case 'outerHTML':
		  elt.outerHTML = swapContent;
		break;
		default: 
		    return;
	      }
	    return
	}
    }

    api.getCurrentFromQueue = () => {
      const q = api.globalQueue 
      const i = api.globalQueueIndex
      return q[i]
    }

    api.getCurrentId = () => {
	const audioElt = api.getAudioElt()
	if (!audioElt) {
	  return ""
	}
	const currentId = audioElt.dataset.audioId
	if (currentId) {
	  return currentId
	}
      return ""
    }

    api.getPlayingId = () => {
      const current = api.getCurrentId()
      return api.isPaused() ? "" : current
    }

    api.skipBackward = () => {
	const audioElt = api.getAudioElt()
	if (!audioElt) return
	if (audioElt.currentTime < 5) {
	  api.globalQueueIndex--
	  api.playerRequest()
	} else {
	  audioElt.currentTime = 0
	} 
    }

    api.skipForward = () => {
	api.globalQueueIndex++ 
	api.playerRequest()
    }

    api.getHxAudioIds = () => {
	const elts = api.getHxAudioElts()
	if (!elts) return
	return elts.map((elt) => {
	    return elt.getAttribute("hx-audio")
	})
    }

    api.newGlobalQueueFromElt = (triggeringElt) => {
	const allElts = api.getHxAudioElts()
	const newQueue = []
	allElts.forEach((elt, index) => {
	    if (elt == triggeringElt) {
		api.globalQueueIndex = index
	    }
	    newQueue.push(elt.getAttribute('hx-audio'))
	})
	api.globalQueue = newQueue
    }


    // could probably append the form data here, but not sure how to add formdata with htmx.ajax
    api.playerRequest = (triggeringElt) => {
	  htmx.ajax('GET', '/player', {
		source: triggeringElt,
		target: "#player",
	  })
    }

    // When playing from an [hx-audio] elt, if the audioId is the same as the currentId, toggle the audio
    // otherwise request a new player
    api.playFromElt = (triggeringElt, event) => {
	const newId = triggeringElt.getAttribute('hx-audio')
	const currentId = api.getCurrentId()
	if (newId === currentId && currentId) {
	  api.toggle()
	} else {
	  api.newGlobalQueueFromElt(triggeringElt)
	  api.playerRequest(triggeringElt)
	}
      }

    api.getNextThree = () => {
      const fullQueue = api.globalQueue
      const index = api.globalQueueIndex
      return fullQueue.slice(index + 1, index + 4) 
    }


      /// Canvases //////////////////////////////////////
    api.addCanvases = (elts) => {
	const currentId = api.getCurrentId()

	elts.forEach((elt) => {
	  const audioId = elt.getAttribute('hx-audio-vis')
	  addVis(elt)
	  if (audioId == currentId) {
	    api.setCurrentCanvas(elt) 
	  }
	  api.canvases.add(elt)
	})
    }

    api.setCurrentCanvas = (elt) => { 
	const currentCanvas = api.getCurrentCanvas()
	if (currentCanvas) { 
	  currentCanvas.onEnd()
	}
	api.currentCanvas = elt
    }

    api.getCurrentCanvas = () => {
	return api.currentCanvas
    }

    api.updateCurrentCanvas = () => {
      const currentId = api.getCurrentId()

      api.canvases.forEach((elt) => {
	const canvasId = elt.getAttribute('hx-audio-vis')
	if (canvasId == currentId) {
	      api.setCurrentCanvas(elt)
	}
      })
    }

    api.removeCanvases = (elts) => {
	  api.canvases = api.canvases.difference(new Set(elts))
    }

    console.log("api \n", api)
    hxAPI = api
    return api
  }
})


	
//////////////////////////////////////////////////////////// other code




  function downloadRegular(elt) {
		const link = document.createElement("a")
		link.href = elt.dataset.audioSrc
		link.download = elt.dataset.audioTitle
		audioDownloadPost(elt.dataset.audioUser, elt.dataset.audioTitle)
		document.body.appendChild(link)
		link.click()
  }


  function downloadStems(elt) {
      const stemsSrcs = JSON.parse(elt.dataset.stemsSrcs)
      var fileUrls  = []
      for (let i = 0; i < stemsSrcs.length; i++) {
	    fileUrls.push(stemsSrcs[i])
      }
      var tempLink = document.createElement("a")
      document.body.appendChild(tempLink)
      downloadLinks(fileUrls)
      function downloadLinks(links) {
	setTimeout(function() {
	  let fileIndex = links.length - 1;
	  fileUrl = links[fileIndex]
	  tempLink.setAttribute("href", fileUrl)
	  tempLink.click()

	  if (fileIndex > -1) {
	    fileUrls.splice(fileIndex, 1);
	  }

	  if (fileUrls.length > 0) { 
	    downloadLinks(fileUrls)
	  } else {
	    document.body.removeChild(tempLink)
	  }
	}, 500)
      }
  }

  async function downloadLink(href) {
		let link = document.createElement("a")
		link.href = href
		link.style.display = 'none'
		link.download = null
		console.log(link)
		link.click()
  }


  function downloadBoth(elt) {
      const stemsSrcs = JSON.parse(elt.dataset.stemsSrcs)
      var fileUrls  = []
      fileUrls.push(elt.dataset.audioSrc)
	for (let i = 0; i < stemsSrcs.length; i++) {
	      fileUrls.push(stemsSrcs[i])
	}
      var tempLink = document.createElement("a")
      document.body.appendChild(tempLink)
      downloadLinks(fileUrls)
      function downloadLinks(links) {
	setTimeout(function() {
	  let fileIndex = links.length - 1;
	  fileUrl = links[fileIndex]
	  tempLink.setAttribute("href", fileUrl)
	  tempLink.click()

	  if (fileIndex > -1) {
	    fileUrls.splice(fileIndex, 1);
	  }

	  if (fileUrls.length > 0) { 
	    downloadLinks(fileUrls)
	  } else {
	    document.body.removeChild(tempLink)
	  }
	}, 500)
      }
  }


    function downloadItem(elt, downloadType) {
	    switch (downloadType) {
	      case "regular": 
		downloadRegular(elt)
		break; 
	      case "stems-only":
		downloadStems(elt)
		break;
	      case "both": 
		downloadBoth(elt)
		break; 
	      default:
		return;
	    }
	  audioDownloadPost(elt.dataset.audioId)
    }


  async function audioDownloadPost(audioId) {
  const url = `/downloads/${audioId}`

    try {
      const response = await fetch(url, {
	method: 'POST',
      })
      if (!response.ok) {
	console.error(response)
      } 
      } catch(e) {
	console.error(e)
      }
    }

//
// function ModalFocusOut(elt, event) {
//   if (elt.contains(event.explicitOriginalTarget)) return
//     elt.innerHTML = ''
//     const currentUrl = window.location.href
//     const newUrl = new URL(currentUrl)
//     newUrl.searchParams.delete("modal")
//     window.history.replaceState({}, null, newUrl)
// }

htmx.on("exit-modal", function() {
    htmx.find("#modal").innerHTML = ""
})


// Canvas / audio vis stuff ///////////////////////////////////////

function fillCanvasBar(elt, i, color, data) {
      const max = Math.max(...(data))
      const barwidth = (elt.scrollWidth / data.length) / 2
      const bargap = (elt.scrollWidth / data.length) / 2
      const rect = elt.getBoundingClientRect();
      const ctx = elt.getContext("2d")
      ctx.fillStyle = color;
      const barheight = rect.height * (data[i] / max)
      const x = i * (barwidth + bargap);
      const y = (rect.height / 2) - (barheight) / 2
      ctx.clearRect(x - bargap / 2, 0, barwidth + bargap, rect.height)
      ctx.fillRect(x, y, barwidth, barheight)
  }

  function addVis(canvas) {
    const ctx = canvas.getContext("2d")
    const data = JSON.parse(canvas.getAttribute("audio-vis-data"))
    const barwidth = (canvas.scrollWidth / data.length) / 2
    const bargap = (canvas.scrollWidth / data.length) / 2 

    const dpr = window.devicePixelRatio;
    const rect = canvas.getBoundingClientRect();
    const max = Math.max(...(data))

    // Set the "actual" size of the canvas
    canvas.width = rect.width * dpr;
    canvas.height = rect.height * dpr;

    // Scale the context to ensure correct drawing operations
    ctx.scale(dpr, dpr);

    function fillBar(i, color) {
      ctx.fillStyle = color;
      const barheight = rect.height * (data[i] / max)
      const x = i * (barwidth + bargap);
      const y = (rect.height / 2) - (barheight) / 2
      ctx.clearRect(x - bargap / 2, 0, barwidth + bargap, rect.height)
      ctx.fillRect(x, y, barwidth, barheight)
    }

    function fillBars(color, divisor) {
      ctx.fillStyle = color
      for (let i = 0; i < data.length; i++) {
        barheight = rect.height * (data[i] / max) / divisor
        ctx.fillRect(i * (barwidth + bargap), (rect.height / 2) - barheight / 2, barwidth, barheight)
      }
    }
    fillBars("white", 1);

    const statusArr = new Array(data.length).fill(false)
    
  canvas.timeUpdate = (progress) => {
      const index = data.length * progress

      for (let i = 0; i < statusArr.length; i++) {
        if (statusArr[i] == false && i <= index) {
          fillBar(i, "#898b8f")
          statusArr[i] = true;
        } else if (statusArr[i] == true && i > index) {
          fillBar(i, "white")
          statusArr[i] = false;
        }
      }
    }

  canvas.onEnd = () => {
      for (let i = 0; i < data.length; i++) {
        fillBar(i, "white")
      }
    }
  }

