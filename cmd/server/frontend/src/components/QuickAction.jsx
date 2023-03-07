export default function QuickAction({ children, onClick, description, selected, className = '' }) {
  return (
    <div className={`cursor-pointer ${className}`}>
      <div
        onClick={onClick}
        className={`rounded-lg ${selected ? 'bg-blue-400' : 'bg-neutral-700'} py-2.5 px-2.5 mx-2.5 my-2 hover:bg-blue-400 duration-75 shadow-sm flex justify-center`}>
        {children}
      </div>
      <p className="text-xs">
        {description}
      </p>
    </div>
  )
}