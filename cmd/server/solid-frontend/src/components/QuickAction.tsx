import { Component, JSX } from "solid-js"

type Props = {
  children: JSX.Element
  onClick: VoidFunction
  description?: string
  selected?: boolean
  class?: string
}

const QuickAction: Component<Props> = (props) => {
  return (
    <div class={`cursor-pointer ${props.class}`}>
      <div
        onClick={props.onClick}
        class={`rounded-lg ${props.selected ? 'bg-blue-400' : 'bg-neutral-700'} py-2.5 px-2.5 mx-2.5 my-2 hover:bg-blue-400 duration-75 shadow-sm flex justify-center`}>
        {props.children}
      </div>
      <p class="text-xs">
        {props.description}
      </p>
    </div>
  )
}

export default QuickAction