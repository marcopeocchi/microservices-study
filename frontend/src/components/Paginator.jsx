import { Button } from "./Button";

export default function Paginator({ pageNumber, currentPage, onClick, onChange }) {
  return (
    <div className='flex flex-row justify-center items-center mb-6'>
      {
        new Array(pageNumber).fill().map((_, index) => (
          <Button
            selected={currentPage === index + 1}
            key={index}
            onClick={() => {
              onClick(index + 1)
              onChange()
            }}
          >
            {index + 1}
          </Button>
        ))
      }
    </div>
  )
}