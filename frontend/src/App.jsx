import { useEffect, useRef, useState } from 'react'
import {
  CgCalendar,
  CgEditFade,
  CgInfo,
  CgList,
  CgLogOut,
  CgRedo,
  CgSearch,
  CgSortAz
} from 'react-icons/cg'
import { useNavigate } from "react-router-dom"
import { debounceTime, fromEvent, map } from 'rxjs'
import ListTile from './components/ListTile'
import Logo from './components/Logo'
import Paginator from './components/Paginator'
import QuickAction from './components/QuickAction'
import { SearchModal } from './components/SearchModal'
import Spinner from './components/Spinner'
import Thumbnail from './components/Thumbnail'
import { getHost, isOrderedByDate, isOrderedByName } from './utils'

function App() {
  const [fetchMode, setFetchMode] = useState(localStorage.getItem("fetch-mode")) || "date"
  const [searchfilter, setSearchFilter] = useState("")
  const [showSearch, setShowSearch] = useState(false)
  const [hide, setHide] = useState(localStorage.getItem("hide") === "true")
  const [loading, setLoading] = useState(true)
  const [pages, setPages] = useState(0)
  const [page, setPage] = useState(1)
  const [list, setList] = useState([])
  const [listView, setListView] = useState(localStorage.getItem("listView") === "true")
  const navigate = useNavigate()

  const main = useRef(null)
  const search = useRef(null)

  useEffect(() => {
    loadData()
  }, [fetchMode, searchfilter, page])

  useEffect(() => {
    const $filter = fromEvent(search.current, 'keyup')
      .pipe(
        debounceTime(500),
        map(e => e.target.value),
      )
      .subscribe(setSearchFilter)
    return () => $filter.unsubscribe()
  }, [])

  const loadData = async () => {
    const res = await fetch(`${getHost()}/list?fetchBy=${fetchMode}&filter=${searchfilter}&page=${page}&pageSize=${56}`)
    if (res.redirected || res.status != 200) {
      navigate('/login')
    }
    const data = await res.json()
    setList(data.list)
    setPages(data.pages)
    setLoading(false)
  }

  const toggleHideThumbnails = () => {
    localStorage.setItem("hide", !hide)
    setHide(state => !state)
  }

  const toggleListView = () => {
    localStorage.setItem("listView", !listView)
    setListView(state => !state)
  }

  return (
    <div className='flex flex-row'>
      <nav className='xl:basis-[3.25%] basis-[5%] shrink-0 bg-neutral-800 flex-col justify-center text-center p-1'>
        <QuickAction
          className='mt-8'
          selected={isOrderedByDate()}
          description={"Date"}
          onClick={() => {
            localStorage.setItem("fetch-mode", "date")
            setFetchMode("date")
          }}>
          <CgCalendar />
        </QuickAction>
        <QuickAction
          selected={isOrderedByName()}
          description="Name"
          onClick={() => {
            localStorage.setItem("fetch-mode", "name")
            setFetchMode("name")
          }}>
          <CgSortAz />
        </QuickAction>
        <QuickAction
          className='mt-8'
          description="List"
          onClick={toggleListView}
          selected={listView}
        >
          <CgList />
        </QuickAction>
        <QuickAction
          description="Hide"
          onClick={toggleHideThumbnails}
          selected={hide}
        >
          <CgEditFade />
        </QuickAction>
        <QuickAction
          description="Reload"
          onClick={() => loadData()}
        >
          <CgRedo />
        </QuickAction>
        <QuickAction
          description="Search"
          onClick={() => setShowSearch(state => !state)}
          selected={showSearch}
        >
          <CgSearch />
        </QuickAction>
        <QuickAction
          className="mt-8"
          description="Help"
          onClick={() => navigate('/help')}
        >
          <CgInfo />
        </QuickAction>
        <QuickAction
          description="Logout"
          onClick={() => navigate('/login')}
        >
          <CgLogOut />
        </QuickAction>
      </nav>
      <SearchModal inputRef={search} show={showSearch} hideCallback={() => {
        setShowSearch(false)
      }} />
      <main className='xl:basis-[96.75%] basis-[95%] overflow-y-scroll h-screen' ref={main}>
        {loading &&
          <div className='flex justify-center items-center h-full'>
            <Spinner />
          </div>
        }
        <div className='px-6 mx-auto'>
          <Logo />
          {!listView &&
            <div className='grid grid-cols-2 2xl:grid-cols-7 xl:grid-cols-5 md:grid-cols-4 gap-2 pt-8 min-h-screen'>
              {(list ?? []).filter(entry => entry.name !== "").map((entry) => (
                <Thumbnail entry={entry} isHidden={hide} key={entry.name} />
              ))}
            </div>
          }
          {listView &&
            <div className='min-h-screen mt-6'>
              <div className='grid grid-cols-1 gap-2'>
                {(list ?? []).filter(entry => entry.name !== "").map((entry) => (
                  <ListTile entry={entry} />
                ))}
              </div>
            </div>
          }
          {pages !== 0 && <Paginator pageNumber={pages} currentPage={page} onClick={setPage} onChange={() => {
            main.current.focus()
            main.current.scrollTo(0, 0)
          }}
          />}
        </div>
      </main>
    </div >
  )
}

export default App
