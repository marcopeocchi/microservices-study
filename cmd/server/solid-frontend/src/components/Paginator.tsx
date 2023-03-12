import { Component, For } from "solid-js";
import Button from "./Button";

type Props = {
  pageNumber: number
  currentPage: number
  onClick: Function
  onChange: Function
}

const Paginator: Component<Props> = ({ pageNumber, currentPage, onClick, onChange }) => {
  return (
    <div class='flex flex-row justify-center items-center mb-6'>
      <For each={new Array(pageNumber).fill(0)}>{(_, index) =>
        <Button
          selected={currentPage === index() + 1}
          onClick={() => {
            onClick(index() + 1)
            onChange()
          }}
        >
          {index() + 1}
        </Button>
      }</For>
    </div>
  )
}

export default Paginator