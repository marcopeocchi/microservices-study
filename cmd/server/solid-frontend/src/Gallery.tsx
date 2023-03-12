import { useLocation, useNavigate, useRouteData } from "@solidjs/router"
import clsx from "clsx"
import { CgArrowLeft, CgArrowRight, CgBackspace, CgExpand, CgSandClock, CgViewSplit } from "solid-icons/cg"
import { Component, createEffect, createResource, createSignal, Show } from "solid-js"
import MediaPlayer from "./components/MediaPlayer"
import QuickAction from "./components/QuickAction"
import Spinner from "./components/Spinner"
import { getHostOverlay, getHostStatic } from "./utils/url"

type GalleryResponse = {
  avif: string[]
  avifAvailable: boolean
  webp: string[]
  webpAvailable: boolean
  source: string[]
  cached: boolean
}

async function fetcher(path: string) {
  const res = await fetch(`${getHostOverlay()}/gallery?dir=${path}`)
  if (res.redirected || res.status != 200) {
    throw new Error()
  }
  const data: GalleryResponse = await res.json()

  if (data.avifAvailable) {
    return data.avif
  }
  if (data.webpAvailable) {
    return data.webp
  }
  return data.source
}

const Gallery: Component = () => {
  let main: HTMLDivElement

  let slideshowTimer: number

  const [_, path] = useLocation().pathname.split('/gallery/')
  const [data] = createResource(path, fetcher)

  const navigate = useNavigate()
  const [index, setIndex] = createSignal(0)

  const [mediaLoading, setMediaLoading] = createSignal(true)
  const [autoEnabled, setAutoEnabled] = createSignal(false)
  const [showCounter, setShowCounter] = createSignal(false)
  const [fullscreen, setFullscreen] = createSignal(false)
  const [splitView, setSplitView] = createSignal(false)

  const handleNext = () => {
    setShowCounter(true)
    setMediaLoading(true)
    setIndex((index() + 1) % data().length)
  }

  const handlePrev = () => {
    setShowCounter(true)
    setMediaLoading(true)
    setIndex((index() <= 0 ? data().length - 1 : index() - 1) % data().length)
  }

  createEffect(() => {
    if (data.error) {
      navigate('/login')
    }
  })

  createEffect(() => {
    if (showCounter()) {
      setTimeout(() => setShowCounter(false), 750)
    }
  })

  createEffect(async () => {
    if (autoEnabled()) {
      await main.requestFullscreen()
      slideshowTimer = setInterval(() => handleNext(), 5000)
      return
    }
    if (document.fullscreenElement) {
      await document.exitFullscreen()
    }
    clearInterval(slideshowTimer)
  })

  createEffect(() => {
    document.onkeyup = (e: KeyboardEvent) => {
      switch (e.key) {
        case 'ArrowLeft':
          handlePrev()
          break
        case 'ArrowRight':
          handleNext()
          break
        case 'Backspace':
          navigate('/')
          break
        case 'v':
          setSplitView(state => !state)
          break
        case 's':
          setFullscreen(state => !state)
          break
        case 'p':
          setAutoEnabled(false)
          break
        default:
          break
      }
    }
  })

  return (
    <div class="bg-neutral-900 min-h-screen" ref={main}>
      <nav class="fixed w-14 h-full bg-transparent hover:bg-neutral-800 flex flex-col items-center justify-between px-8 duration-100 opacity-90 text-center group">
        <div class="flex flex-col items-center mt-8">
          <QuickAction
            onClick={() => navigate("/")}
            description="Back"
            class="hidden group-hover:group-hover:block"
          >
            <CgBackspace />
          </QuickAction>
          <QuickAction
            onClick={() => setSplitView(state => !state)}
            description="Split"
            class="hidden group-hover:group-hover:block"
          >
            <CgViewSplit />
          </QuickAction>
          <QuickAction
            onClick={() => setFullscreen(state => !state)}
            description="Span-H"
            class="hidden group-hover:group-hover:block"
          >
            <CgExpand />
          </QuickAction>
          <QuickAction
            selected={autoEnabled()}
            onClick={() => setAutoEnabled(state => !state)}
            description="Auto"
            class="hidden group-hover:group-hover:block"
          >
            <CgSandClock />
          </QuickAction>
        </div>
        <div>
          <QuickAction
            onClick={handleNext}
            class="hidden group-hover:group-hover:block"
          >
            <CgArrowLeft />
          </QuickAction>
          <QuickAction
            onClick={handlePrev}
            class="hidden group-hover:group-hover:block"
          >
            <CgArrowRight />
          </QuickAction>
        </div>
        <div></div>
        <div></div>
        <div class="hidden group-hover:group-hover:block text-sm mb-3">
          Fuu v1.2
        </div>
      </nav>
      <Show when={!data.loading} fallback={<Spinner />}>
        <div class="flex items-center justify-center">
          <Show when={mediaLoading()}>
            <div class="absolute top-1/2 left-1/2 transform -translate-x-1/2 -translate-y-1/2 z-50">
              <Spinner />
            </div>
          </Show>
          <Show when={!splitView()}>
            <MediaPlayer
              onLoad={() => setMediaLoading(false)}
              loading={mediaLoading()}
              fullscreen={fullscreen()}
              src={`${getHostStatic()}/${path}/${data().at(index())}`}
            />
          </Show>
          <Show when={splitView()}>
            <MediaPlayer
              onLoad={() => setMediaLoading(false)}
              loading={mediaLoading()}
              fullscreen={fullscreen()}
              src={`${getHostStatic()}/${path}/${data().at(index())}`}
            />
            <MediaPlayer
              loading={mediaLoading()}
              onLoad={() => setMediaLoading(false)}
              fullscreen={fullscreen()}
              src={`${getHostStatic()}/${path}/${data().at((index() + 1) % data().length)}`}
            />
          </Show>
        </div>
        <footer class={clsx(
          showCounter() ? 'bg-neutral-800' : 'bg-transparent',
          'fixed bottom-8 left-1/2 -translate-x-1/2 w-40 h-40 rounded-lg flex flex-col items-center justify-center px-3 duration-150 opacity-90 text-center group'
        )}>
          <div class="flex justify-center">
            <div class={clsx(
              showCounter() ? 'block' : 'hidden',
              'font-semibold text-5xl text-blue-400'
            )}>
              {Number(index()) + 1} / {data().length}
            </div>
          </div>
        </footer>
      </Show>
    </div>
  )
}

export default Gallery