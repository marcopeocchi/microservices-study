import { useNavigate } from '@solidjs/router'
import {
  CgCalendar,
  CgShield,
  CgInfo,
  CgList,
  CgLogOut,
  CgSearch,
  CgSortAz
} from 'solid-icons/cg'
import { Component, createEffect, createResource, createSignal, For, Show } from 'solid-js'
import ListTile from './components/ListTile'
import Logo from './components/Logo'
import Paginator from './components/Paginator'
import QuickAction from './components/QuickAction'
import SearchModal from './components/SearchModal'
import Spinner from './components/Spinner'
import Thumbnail from './components/Thumbnail'
import { useLogout } from './hooks/useLogout'
import { isOrderedByDate, isOrderedByName } from './utils/files'
import { composeQuery, getHostOverlay } from './utils/url'


const App: Component = () => {
  const [fetchMode, setFetchMode] = createSignal(localStorage.getItem("fetch-mode") || "date")
  const [filter, setFilter] = createSignal("")
  const [showSearch, setShowSearch] = createSignal(false)
  const [hide, setHide] = createSignal(localStorage.getItem("hide") === "true")
  const [page, setPage] = createSignal(Number(window.location.hash.split('-').at(1)) || 1)
  const [listView, setListView] = createSignal(localStorage.getItem("listView") === "true")

  const navigate = useNavigate()
  const logout = useLogout()

  const derivedSignal = () => { return { page: page(), fetchMode: fetchMode(), filter: filter() } }

  const fetcher = async (signal: {
    page: number,
    fetchMode: string,
    filter: string,
  }) => {
    const res = await fetch(`${getHostOverlay()}/${composeQuery(signal.fetchMode, signal.filter, signal.page, 49)}`)
    if (!res.ok) {
      navigate('/login')
      throw new Error(`Error: ${res.status}`)
    }
    await new Promise(s => setTimeout(s, 250))
    const data: Paginated<Directory> = await res.json()
    return data
  }

  const [data] = createResource(derivedSignal, fetcher)


  let main: HTMLElement

  const toggleHideThumbnails = () => {
    localStorage.setItem("hide", String(!hide()))
    setHide(state => !state)
  }

  const toggleListView = () => {
    localStorage.setItem("listView", String(!listView()))
    setListView(state => !state)
  }

  return (
    <div class='flex flex-row'>
      <nav class='xl:basis-[3.25%] basis-[5%] shrink-0 bg-neutral-800 flex-col justify-center text-center p-1'>
        <QuickAction
          class='mt-8'
          selected={isOrderedByDate(fetchMode())}
          description={"Date"}
          onClick={() => {
            localStorage.setItem("fetch-mode", "date")
            setFetchMode("date")
          }}>
          <CgCalendar />
        </QuickAction>
        <QuickAction
          selected={isOrderedByName(fetchMode())}
          description="Name"
          onClick={() => {
            localStorage.setItem("fetch-mode", "name")
            setFetchMode("name")
          }}>
          <CgSortAz />
        </QuickAction>
        <QuickAction
          class='mt-8'
          description="List"
          onClick={toggleListView}
          selected={listView()}
        >
          <CgList />
        </QuickAction>
        <QuickAction
          description="Hide"
          onClick={toggleHideThumbnails}
          selected={hide()}
        >
          <CgShield />
        </QuickAction>
        <QuickAction
          description="Search"
          onClick={() => setShowSearch(state => !state)}
          selected={showSearch()}
        >
          <CgSearch />
        </QuickAction>
        <QuickAction
          class="mt-8"
          description="Help"
          onClick={() => navigate('/help')}
        >
          <CgInfo />
        </QuickAction>
        <QuickAction
          description="Logout"
          onClick={() => {
            logout()
            navigate('/login')
          }}
        >
          <CgLogOut />
        </QuickAction>
      </nav>
      <SearchModal
        show={showSearch()}
        onSearch={(e) => setFilter(e)}
        hideCallback={() => setShowSearch(false)}
      />
      <main class='xl:basis-[96.75%] basis-[95%] overflow-y-scroll h-screen' ref={main}>
        <Show when={!data.loading} fallback={
          <div class='flex justify-center items-center h-full'>
            <Spinner />
          </div>
        }>
          <div class='px-6 mx-auto'>
            <Logo />
            <Show when={!listView()}>
              <div class='grid grid-cols-1 2xl:grid-cols-7 xl:grid-cols-5 md:grid-cols-4 sm:grid-cols-2 gap-2 pt-8 min-h-screen'>
                <For each={data().list.filter(entry => entry.name !== "")}>{(entry) =>
                  <Thumbnail entry={entry} isHidden={hide()} />
                }</For>
              </div>
            </Show>
            <Show when={listView()}>
              <div class='min-h-screen my-6'>
                <div class='grid grid-cols-1 gap-2'>
                  <For each={data().list.filter(entry => entry.name !== "")}>{(entry) =>
                    <ListTile entry={entry} />
                  }</For>
                </div>
              </div>
            </Show>
            <Show when={data().pages !== 0}>
              <Paginator
                pageNumber={data().pages}
                currentPage={data().page}
                onClick={setPage}
                onChange={() => {
                  main.focus()
                  main.scrollTo(0, 0)
                }}
              />
            </Show>
          </div>
        </Show>
      </main>
    </div >
  )
}

export default App
