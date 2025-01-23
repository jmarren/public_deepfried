let uploadApi


htmx.defineExtension('upload', {
 onEvent: function(name, event) {

    /* -------------- afterProcessNode ------------*/
     if (name === 'htmx:afterProcessNode') {
      const node = event.target || event.detail.elt;

      if (node.hasAttribute('hx-file-upload')) {
	    htmx.on(node, 'input', getFileInputChangedListener(node))
       }

      if (node.hasAttribute('hx-upload-drop-zone')) {
	    htmx.on(node, 'drop', getDropListener(node))
       }
     }



      /*--------------- configRequest -----------*/
       if (name === "htmx:configRequest") {
          const elt = event.target || event.detail.elt;
          
          if (elt.hasAttribute('hx-file-upload')) {
              event.detail.headers['X-upload-id'] = elt.getAttribute('hx-file-upload-id')
              event.detail.headers['X-file-name'] = elt.files[0].name
              event.detail.formData = null
          }
       }


 },
  init: function(api) {
    uploadApi = api
    return api
  }
})




// When an upload input changes, create a file object 
// and send it to the worker
function getFileInputChangedListener(elt) {
      return function(evt) {
	const uploadId = elt.getAttribute('hx-file-upload-id')
	const filename = elt.files[0].name
	fileObj = {
		filename: filename,
		totalSize: elt.files[0].size,
		id: uploadId,
		}
	UploadFile(elt)
      }
}


function ShowArtwork(elt) {
  const file = elt.files[0]
  const fileReader = new FileReader()
  fileReader.readAsDataURL(file)
  const imageElt = htmx.find("#artwork-display")
  htmx.on(fileReader, "load", () => {
    imageElt.src = fileReader.result
  })
  const fileNameInput = htmx.find("input[name=artwork_file_name]")
  fileNameInput.value = file.name

}


function UploadFile(elt) {
  const file = elt.files[0]
  const uploadType = elt.getAttribute('hx-file-upload-type')
  const uploadId = elt.getAttribute('hx-file-upload-id')
  const fileObj = { 
	file: file,
	type: uploadType,
	id: uploadId,
   }

 const worker = new Worker("/workers/upload-worker.js")

  worker.postMessage(fileObj)
  worker.addEventListener('message', (e) => {
	UploadWorkerMsgListener(e)
  })
}


function UploadWorkerMsgListener(e) {
    const filename = e.data[0]
    const progress = e.data[1]
    const uploadType = e.data[2]
    const msg = e.data[3]

    if (uploadType == "audio") {
	HandleMainAudioProgress(progress, msg)
    }
}


function HandleMainAudioProgress(progress, msg) {
    const uploadProgressElt = htmx.find("[hx-upload-progress]")
    uploadProgressElt.style.background = `linear-gradient(90deg,black 0%, black ${progress * 100}%, white ${progress *  100}%, white ${100}%)`

    // is msg is complete, add checkmark
    const completionElt = htmx.find("[hx-upload-status]")
    if (msg == "complete") {
	    completionElt.innerHTML =  `
		    <svg xmlns="http://www.w3.org/2000/svg" height="24px" viewBox="0 -960 960 960" width="24px" fill="green"><path d="M382-240 154-468l57-57 171 171 367-367 57 57-424 424Z"/></svg>
	      `
    }
    // if msg has an error, add red X
    if (msg.includes('error')) {
	    completionElt.innerHTML = `
<svg xmlns="http://www.w3.org/2000/svg" height="24px" viewBox="0 -960 960 960" width="24px" fill="red"><path d="m256-200-56-56 224-224-224-224 56-56 224 224 224-224 56 56-224 224 224 224-56 56-224-224-224 224Z"/></svg> upload failed
` 
    }
}




// /* From upload.js */ //
 function getDropListener(elt) {
  const targetSelector = elt.getAttribute('hx-upload-drop-zone')
  return function(e) {
        const targetElt = htmx.find(targetSelector)
	targetElt.files = e.dataTransfer.files
  }
 }




function setUploadErr(msg) {
      const errorElt = htmx.find("#upload-error")
      errorElt.innerHTML = msg
}


/* upload_form.js */
function AddTagButton(elt) {
  const container = htmx.closest(elt, ".input-container")
  const numTagInputs = htmx.findAll(container, "input").length
  if (numTagInputs >= 3) return
  const newTagContainer = document.createElement("div")
  newTagContainer.className="tag-input-container"
  const newInput = document.createElement("input")
  newInput.name=`tag-${numTagInputs + 1}`
  newInput.className="tag-input input-field"
  newInput.placeholder = "jazz"
  newInput.type = "text"
  newTagContainer.appendChild(newInput)
  const newRemoveButton = document.createElement("button")
  newRemoveButton.onclick = (e) => e.target.closest(".tag-input-container").remove()
  newRemoveButton.className="remove-tag-button"
  const removeIcon = document.createElement("div")
  removeIcon.className="remove-icon"
  newRemoveButton.appendChild(removeIcon)
  newTagContainer.appendChild(newRemoveButton)
  container.insertBefore(newTagContainer, elt)
  htmx.process(container)
}

// when the input of the input[name=stems] element changes,
// addStemFile is invoked.
// It will get the names of the stem files and add them to the input[name=stem_file_names] input
// so they can be posted when the upload is confirmed
function addStemFile(elt) {
  console.log("files: ", elt.files)
  const files = elt.files
  const stemsSection = htmx.find("#stems-input-section")
  const container = htmx.find(stemsSection, ".input-container")
  const addButton = htmx.find(stemsSection, "#add-stems-button")

  const stemFileNamesInput = htmx.find("input[name=stem_file_names]") 
  const stemFileNames = JSON.parse(stemFileNamesInput.value)
  console.log(stemFileNames)

  for (let i = 0; i < files.length; i++) {
    const fileElt = document.createElement('div');
    fileElt.innerHTML = files[i].name
    stemFileNames.push(files[i].name)
    fileElt.className = "stem-file"
    container.insertBefore(fileElt, addButton)
    InitUpload(files[i], "stems")
  }
  const newValue = JSON.stringify(stemFileNames)
  console.log(newValue)
  stemFileNamesInput.value = newValue
}



