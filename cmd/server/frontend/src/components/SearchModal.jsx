import clsx from "clsx";
import { useEffect } from "react";
import { CgSearch } from "react-icons/cg";

export function SearchModal({ show, inputRef, hideCallback }) {
  useEffect(() => {
    if (show) {
      inputRef.current.focus()
    }
  }, [inputRef, show])

  return (
    <>
      {show && <div className="fixed w-full min-h-screen bg-neutral-900/50" onClick={hideCallback} />}
      <div className={clsx(
        show ? 'block' : 'hidden',
        'fixed top-1/2 left-1/2 -translate-y-1/2 -translate-x-1/2 bg-neutral-800/80 h-1/4 w-96 rounded-md border-2 border-neutral-500 backdrop-blur-md'
      )}>
        <div className='flex justify-center items-center h-full text-center'>
          <div className='flex flex-col items-center'>
            <div className='font-bold text-neutral-100 text-4xl mb-8 flex items-end'>
              <span>
                Search
              </span>
              <CgSearch />
            </div>
            <input
              type="text"
              className='bg-neutral-800 text-neutral-300 w-80 px-2.5 py-2 rounded border-2 placeholder:text-neutral-500 focus:outline-0 border-blue-400 appearance-none'
              placeholder='Filter album'
              ref={inputRef}
            />
          </div>
        </div>
      </div>
    </>
  )
}