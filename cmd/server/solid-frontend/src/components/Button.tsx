import { Component, JSX } from "solid-js"

type Props = {
  children: JSX.Element
  onClick: VoidFunction
  class?: string
  selected?: boolean
}

const Button: Component<Props> = (props) => {
  return (
    <div
      onClick={props.onClick}
      class={`${props.class} font-semibold group-hover:group-hover:block text-center hover:bg-blue-400 duration-100 ${props.selected ? 'bg-blue-400' : 'bg-neutral-700'} px-2.5 py-1.5 rounded cursor-pointer mr-2`}>
      {props.children}
    </div>
  )
}

export default Button