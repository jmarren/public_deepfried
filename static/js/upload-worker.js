

onmessage = (e) => {
  console.log("upload-worker.js received message")
  const files = e.data
  console.log(files)
  const filesToUpload = []
  for (let i = 0; i < files.length; i++) {
    filesToUpload.push(new FileToUpload(files[i], i))
  }
  const upload = new Upload(filesToUpload)
  upload.uploadFiles()
}

class Upload {
  constructor(filesToUpload) {
    this.filesToUpload = filesToUpload
  }
  uploadFiles() {
    for (let i = 0; i < this.filesToUpload.length; i++) {
      this.filesToUpload[i].upload()
    }
  }
}

class FileToUpload {
  constructor(file, id) {
    this.file = file
    this.fileName = file.name
    this.stream = file.stream()
    this.chunkIndex = 0
    this.id = id
    this.writeStream = this.createWriteStream()
    this.progress = 0
  }
  incrementChunkIndex() {
    this.chunkIndex++
  }
  upload() {
    if (this.file.type != "audio/mpeg" && this.file.type != "audio/wav") {
      this.stream.abort()
    } else {
      this.stream.pipeTo(this.writeStream)
    }
  }
  async postChunk(chunk) {
    const index = `${this.chunkIndex}`
    const fileId = `${this.id}`
    const fileName = `${this.fileName}`
    const fileSize = `${this.file.size}`
    try {

      const chunkSize = `${chunk.byteLength}`

      const response = await fetch('http://localhost:8080/upload', {
        method: 'POST',
        body: chunk,
        headers: {
          'Content-Type': 'audio/mpeg',
          'X-index': index,
          'X-total-size': fileSize,
          'X-file-name': fileName,
          'X-chunk-size': chunkSize,
        }
      })
      console.log("response: ", response)

      if (!response.ok) {
        throw new Error("failed to upload chunk");
      }
    } catch (e) {
      console.error(e)
    } finally {

      this.chunkIndex++
      this.progress = this.progress + chunk.byteLength / fileSize;
      postMessage([this.chunkIndex, this.progress])
    }
  }
  createWriteStream() {
    const queuingStrategy = new CountQueuingStrategy({ highWaterMark: 1 })

    const underlyingSink = {
      start: () => {
        console.log("stream start")
      },
      write: (chunk) => {
        return this.postChunk(chunk)
      },
      close: () => {
        console.log("closing stream")
      },
      abort: (err) => {
        console.error(err)
      },
    }
    const writableStream = new WritableStream(underlyingSink, queuingStrategy)
    return writableStream
  }

}





