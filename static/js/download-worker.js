
onmessage = (e) => {
  console.log("working")
  // const files = e.data
  // console.log(files)
  // const filesToUpload = []
  // for (let i = 0; i < files.length; i++) {
  //   filesToUpload.push(new FileToUpload(files[i], i))
  // }
  // const upload = new Upload(filesToUpload)
  // upload.uploadFiles()
  const baseUrl = e.data[0]
  const numChunks = e.data[1]
  const download = new Download(baseUrl, numChunks)
  download.downloadAll()
}


class Download {
  constructor(baseUrl, numChunks) {
    this.baseUrl = baseUrl
    this.numChunks = numChunks
  }
  downloadAll() {
    for (let i = 1; i < this.numChunks; i++) {
      this.downloadChunk(i)
    }
  }
  async downloadChunk(num) {
    try {
      const chunkUrl = `${this.baseUrl}/chunk-${num}.mp3`
      const response = await fetch(chunkUrl)
      if (!response) {
        throw new Error(`Reponse Status: ${response.status}`)
      }
      // console.log(response)
      const chunk = response.body
      this.processAndPost(chunk)
    } catch (err) {
      console.log(err)
    }
  }
  async processAndPost(chunk) {
    // const audioCtx = new AudioContext()
    console.log(chunk)
    const reader = chunk.getReader()
    try {
      const readData = await reader.read()
      console.log("worker read data: ", readData)
      const arrayBuff = readData.value.buffer
      console.log("arrayBuff worker: ", arrayBuff)
      postMessage(arrayBuff)

    } catch (e) {
      console.error(e)
    }
    // const array = data.value
    // const buffer = array.buffer.slice(array.byteOffset, array.byteLength + array.byteOffset)
    // console.log("buffer: ", buffer)
  }
}
