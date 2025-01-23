/* Moved to upload.js for the time being

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

async function InitUpload(file, type) {
  const url = `/upload/new/${type}`
  formData = new FormData()
  formData.append("filename", file.name)
  formData.append("totalSize", file.size)

  const body = {
    'filename': file.name,
    'totalSize': file.size,
  }

  if (type == "stems") {
    formData.append("mainfile", htmx.find("input[name=audio_file_name]").value)
    // body.mainfile = htmx.find("input[name=audio_file_name]").value 
  }
  try {
    const response = await fetch(url, {
      method: 'POST',
      body: formData,
    })
    if (!response.ok) {
      throw new Error("failed init upload");
    }
    const res = await response.text()
    message = res
    return true

  } catch (e) {
    console.error(e)
    return false
  } finally {
    // const worker = new Worker("/static/assets/js/upload-worker.js")
    const worker = new Worker("https://dvnkl2og6j3o0.cloudfront.net/public/js/upload-worker.js")
    window.workers.push(worker)
    let workerIndex = window.workers.length - 1
    const fileObj = {
      file: file,
      type: type,
    } 
    window.workers[workerIndex].postMessage(fileObj)
  }
}

function HandleArtworkSelected(elt) {
  const file = elt.files[0]
  const fileReader = new FileReader()
  fileReader.readAsDataURL(file)
  const artworkInputContainer = htmx.closest(elt, ".artwork-input-container")
  const imageElt = htmx.find(artworkInputContainer, ".file-image")
  htmx.on(fileReader, "load", () => {
    imageElt.src = fileReader.result
  })
  const selectedFile = htmx.closest(elt, ".selected-file")
  const artworkNameInput = htmx.find(selectedFile, "input[name=artwork_file_name]")
  artworkNameInput.value = file.name
  InitUpload(file, "artwork")
}


function UploadAudioFile(fileObj) {
  // const worker = new Worker("/static/assets/js/upload-worker.js")
 const worker = new Worker("https://dvnkl2og6j3o0.cloudfront.net/public/js/upload-worker.js")
  worker.postMessage(fileObj)
  worker.addEventListener("message", (e) => {
    const filename = e.data[0]
    const progress = e.data[1]
    const msg = e.data[2]
    const fileElt = document.getElementById(`file_${e.data[0]}`)
    const progressElt = htmx.find(fileElt, ".progress-marker")
    progressElt.style.width = `${e.data[1] * 100}%`

    if (msg == "complete") {
      const parent = progressElt.parentElement
      parent.style.backgroundColor = 'white'
      parent.innerHTML = 'upload complete'
      parent.style.border = ''
      parent.style.fontFamily = 'Inter'
      parent.style.color = '#706f6f'
      parent.className = 'success-message'
    } else if (msg.includes("error")) {
      htmx.find("#upload-error").innerHTML = msg		
      const parent = progressElt.parentElement
      parent.style.backgroundColor = 'white'
      parent.innerHTML = 'error'
      parent.style.border = ''
      parent.style.fontFamily = 'Inter'
      parent.style.color = '#706f6f'
      parent.className = 'success-message'

    }
  })
}

function StartAudioFileUpload() {
  const audioFiles = htmx.find("#file-input").files
  for (let i = 0; i < audioFiles.length; i++) {
    const filename = audioFiles[i].name
    const audioFile = {file: audioFiles[i], type: "audio"}
    UploadAudioFile(audioFile)
  }
}

// window.addEventListener("upload-audio-files", function() {
//   console.log("upload-audio-files triggered")
// })
//
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
*/
