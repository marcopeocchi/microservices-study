import clsx from "clsx"
import { Fragment, useEffect, useRef, useState } from "react"
import { CgArrowLeft, CgArrowRight, CgBackspace, CgExpand, CgSandClock, CgViewSplit } from 'react-icons/cg'
import { useNavigate } from 'react-router-dom'
import MediaPlayer from './components/MediaPlayer'
import QuickAction from './components/QuickAction'
import Spinner from './components/Spinner'
import { getHost, getHostStatic } from "./utils"

export default function Gallery() {
  const [galleryTitle] = useState([window.location.pathname.split("/").slice(-1)] ?? "-")
  const [currentIndex, setCurrentIndex] = useState(0)
  const [showCounter, setShowCounter] = useState(false)
  const [galleryData, setGalleryData] = useState([])
  const [autoEnabled, setAutoEnabled] = useState(false)
  const [fullscreen, setFullScreen] = useState(false)
  const [splitView, setSplitView] = useState(false)
  const [loading, setLoading] = useState(true)

  const navigate = useNavigate()
  const slideshowTimer = useRef()

  const imgLoaded = () => {
    setLoading(false)
    showCounterTimeout()
  }

  const videoLoaded = () => {
    setLoading(false)
    showCounterTimeout()
  }

  const reset = () => {
    setLoading(true)
  }

  const showCounterTimeout = () => {
    setShowCounter(true)
    setTimeout(() => setShowCounter(false), 1000)
  }

  const startSlideShow = async () => {
    await document.getElementById('galleryMain').requestFullscreen()
    setAutoEnabled(state => !state)
    setTimeout(handleNext, 1000)
    slideshowTimer.current = setInterval(handleNext, 5000)
  }

  const stopSlideShow = async () => {
    await document.exitFullscreen()
    if (slideshowTimer.current) {
      clearInterval(slideshowTimer.current)
    }
    setAutoEnabled(state => !state)
  }

  useEffect(() => {
    (async () => {
      const res = await fetch(`${getHost()}/overlay/gallery?dir=${galleryTitle}`)
      if (res.redirected || res.status != 200) {
        navigate('/login')
      }
      const data = await res.json()

      if (data.avifAvailable) {
        setGalleryData(data.avif)
      } else if (data.webp) {
        setGalleryData(data.source)
      } else {
        setGalleryData(data.source)
      }
    })();
  }, [galleryTitle])

  useEffect(() => {
    const handler = (event) => {
      switch (event.key) {
        case 'ArrowLeft':
          handlePrev()
          window.scrollTo({ top: 0, left: 0, behavior: "auto" })
          break
        case 'ArrowRight':
          handleNext()
          window.scrollTo({ top: 0, left: 0, behavior: "auto" })
          break
        case 'v':
          setSplitView(state => !state)
          break
        case 's':
          setFullScreen(state => !state)
          break
        case 'p':
          autoEnabled ? stopSlideShow() : startSlideShow()
          break
        case 'Backspace':
          navigate('/')
          break
        default: break
      }
    }
    document.onkeydown = handler
  }, [currentIndex, galleryData])

  const handleNext = () => {
    if (splitView) {
      setCurrentIndex(state => (state + 2) % galleryData.length)
    } else {
      setCurrentIndex(state => (state + 1) % galleryData.length)
    }
    reset()
  }
  const handlePrev = () => {
    if (splitView) {
      setCurrentIndex(state => (state - 2) % galleryData.length)
    } else {
      setCurrentIndex(state => (state - 1) < 0 ? (galleryData.length - 1) : (state - 1))
    }
    reset()
  }

  const handleMobileTap = (event) => {
    event.clientX > window.innerWidth / 2 ? handlePrev() : handleNext()
  }

  return (
    <main id="galleryMain">
      <nav className="fixed w-14 h-full bg-transparent hover:bg-neutral-800 flex flex-col items-center justify-between px-8 duration-100 opacity-90 text-center group">
        <div className="flex flex-col items-center mt-8">
          <QuickAction
            onClick={() => navigate("/")}
            description="Back"
            className="hidden group-hover:group-hover:block"
          >
            <CgBackspace />
          </QuickAction>
          <QuickAction
            onClick={() => setSplitView(state => !state)}
            description="Split"
            className="hidden group-hover:group-hover:block"
          >
            <CgViewSplit />
          </QuickAction>
          <QuickAction
            onClick={() => setFullScreen(state => !state)}
            description="Span-H"
            className="hidden group-hover:group-hover:block"
          >
            <CgExpand />
          </QuickAction>
          <QuickAction
            selected={autoEnabled}
            onClick={() => autoEnabled ? stopSlideShow() : startSlideShow()}
            description="Auto"
            className="hidden group-hover:group-hover:block"
          >
            <CgSandClock />
          </QuickAction>
        </div>
        <div>
          <QuickAction
            onClick={handleNext}
            className="hidden group-hover:group-hover:block"
          >
            <CgArrowLeft />
          </QuickAction>
          <QuickAction
            onClick={handlePrev}
            className="hidden group-hover:group-hover:block"
          >
            <CgArrowRight />
          </QuickAction>
        </div>
        <div></div>
        <div></div>
        <div className="hidden group-hover:group-hover:block text-sm mb-3">
          Fuu v1.2
        </div>
      </nav>
      <aside className="fixed mt-8 mr-8 right-0 w-28 h-14 rounded-lg bg-transparent hover:bg-neutral-800 flex flex-col items-center justify-center px-2 duration-100 opacity-90 text-center group">
        <div className="hidden group-hover:group-hover:block w-full">
          <select
            value={currentIndex}
            className='rounded px-2.5 py-1.5 bg-neutral-800'
            onChange={(e) => setCurrentIndex(e.target.value)}>
            {
              galleryData.map((_, idx) => (
                <option value={idx} key={idx}>{idx + 1}</option>
              ))
            }
          </select>
        </div>
      </aside>
      <footer className={clsx(
        showCounter ? 'bg-neutral-800' : 'bg-transparent',
        'fixed bottom-8 left-1/2 -translate-x-1/2 w-40 h-40 rounded-lg flex flex-col items-center justify-center px-3 duration-150 opacity-90 text-center group'
      )}>
        <div className="flex justify-center">
          <div className={clsx(
            showCounter ? 'block' : 'hidden',
            'font-semibold text-5xl text-pink-400'
          )}>
            {Number(currentIndex) + 1} / {galleryData.length}
          </div>
        </div>
      </footer>
      <div className="flex justify-center items-center">
        {galleryData.length &&
          !splitView ?
          <MediaPlayer
            onClick={handleMobileTap}
            onImageLoad={imgLoaded}
            onVideoLoad={videoLoaded}
            fullscreen={fullscreen}
            loading={loading}
            src={`${getHostStatic()}/${galleryTitle}/${galleryData.at(currentIndex)}`}
          /> :
          <Fragment>
            <MediaPlayer
              onClick={handleMobileTap}
              onImageLoad={imgLoaded}
              onVideoLoad={videoLoaded}
              fullscreen={fullscreen}
              loading={loading}
              src={`${getHostStatic()}/${galleryTitle}/${galleryData.at(currentIndex + 1)}`}
            />
            <MediaPlayer
              onClick={handleMobileTap}
              onImageLoad={imgLoaded}
              onVideoLoad={videoLoaded}
              onEnded={handleNext}
              fullscreen={fullscreen}
              loading={loading}
              src={`${getHostStatic()}/${galleryTitle}/${galleryData.at(currentIndex)}`}
            />
          </Fragment>
        }
      </div>
      {loading &&
        <div className="absolute top-1/2 left-1/2 transform -translate-x-1/2 -translate-y-1/2">
          <Spinner />
        </div>
      }
    </main>
  )
}