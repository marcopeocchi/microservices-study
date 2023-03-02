export function Button({ children, onClick, className, selected = false }) {
  return (
    <div
      onClick={onClick}
      className={`${className} font-semibold group-hover:group-hover:block text-center hover:bg-blue-400 duration-100 ${selected ? 'bg-blue-400' : 'bg-neutral-700'} px-2.5 py-1.5 rounded cursor-pointer mr-2`}>
      {children}
    </div>
  )
}