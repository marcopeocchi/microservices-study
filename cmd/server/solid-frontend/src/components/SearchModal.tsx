import clsx from "clsx";
import { debounceTime, fromEvent, map } from "rxjs";
import { CgSearch } from "solid-icons/cg";
import { Component, createEffect } from "solid-js";

type Props = {
  show?: boolean
  hideCallback: VoidFunction
  onSearch: (filter: string) => any
}

const SearchModal: Component<Props> = (props) => {
  let input: HTMLInputElement

  createEffect(() => {
    if (props.show) {
      input.focus()
    }
  })

  createEffect(() => {
    const $filter = fromEvent<any>(input, 'keyup')
      .pipe(
        debounceTime(500),
        map(e => e.target.value),
      )
      .subscribe(props.onSearch)
    return () => $filter.unsubscribe()
  })

  return (
    <>
      {props.show && <div class="fixed w-full min-h-screen bg-neutral-900/50" onClick={props.hideCallback} />}
      <div class={clsx(
        props.show ? 'block' : 'hidden',
        'fixed top-1/2 left-1/2 -translate-y-1/2 -translate-x-1/2 bg-neutral-800/80 h-1/4 w-96 rounded-md border-2 border-neutral-500 backdrop-blur-md'
      )}>
        <div class='flex justify-center items-center h-full text-center'>
          <div class='flex flex-col items-center'>
            <div class='font-bold text-neutral-100 text-4xl mb-8 flex items-end'>
              <span>
                Search
              </span>
              <CgSearch />
            </div>
            <input
              type="text"
              class='bg-neutral-800 text-neutral-300 w-80 px-2.5 py-2 rounded border-2 placeholder:text-neutral-500 focus:outline-0 border-blue-400 appearance-none'
              placeholder='Filter album'
              ref={input}
            />
          </div>
        </div>
      </div>
    </>
  )
}

export default SearchModal