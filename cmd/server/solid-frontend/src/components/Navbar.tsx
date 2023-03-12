import { createEffect } from "solid-js"


document.documentElement.classList.add('dark')

export function Navbar(props: { onChange: Function }) {
  return (
    <div class="sticky top-0 z-50 bg-neutral-50 dark:bg-neutral-800 h-16 w-full shadow-sm py-2.5 px-6 flex items-center justify-between">
      <div class="text-3xl font-bold text-transparent bg-clip-text bg-gradient-to-r from-blue-500 to-blue-400">Fuu</div>
      <input
        placeholder="Search album"
        onInput={e => props.onChange(e.currentTarget.value.length > 2 ? e.currentTarget.value : '')}
        class="rounded-xl text-neutral-700 dark:text-neutral-200 bg-white dark:bg-neutral-900/50 w-96 h-full border border-neutral-100 dark:border-neutral-900 px-4 text-center focus:outline-0 focus:ring-2 focus:ring-blue-500"
      />
      <div class="flex flex-col">
        <div class="text-neutral-500 dark:text-neutral-400">fe-solid.beta</div>
      </div>
    </div>
  )
}