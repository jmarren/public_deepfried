function() {

  var hxAPI

htmx.defineExtension('audio', {
 onEvent: function(name, event) {
    // add the event listener in the afterProcessNode event
     if (name === 'htmx:afterProcessNode') {
        const parent = event.target || event.detail.elt;
        const audioNodes = [
          ...parent.hasAttribute("hx-audio") ? [parent] : [],
          ]
        audioNodes.forEach(function(node) {
          // Initialize the node with the `hx-audio` attribute
          init(node)
        }
        }
  
     if (name == "htmx:configRequest") {
        if (event.target.hasAttribute("hx-audio") || event.detail.elt.hasAttribute("hx-audio") {
          console.log("event: \n", event)
        } 
      }

        return
  },

  init: function(api) {
    api.audioQueue = []
    api.audioPlaying = ""

    api.setQueue = (queue) => {
      api.audioQueue = queue
    }
    api.getQueue= () => {
      return api.audioQueue
    }
    api.setPlaying = (newAudio) => {
      api.audioPlaying = newAudio 
    }
    api.pausePlaying = () =>  {
      api.audioPlaying = ""
    }
    api.getPlaying = () => {
      return api.audioPlaying
    }

    console.log("api: \n", api)
    hxAPI = api
    return api

  }})

  function init(node) {
    node.addEventListener('click', 

  })




