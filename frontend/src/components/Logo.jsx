export default function Logo({ hideSubText = false }) {
  return (
    <>
      <h1 className='pt-7 text-center font-extrabold text-transparent text-6xl bg-clip-text bg-gradient-to-r from-blue-600 to-sky-500'>
        Fuu
      </h1>
      {!hideSubText && <div className='text-center text-xl text-neutral-100'>
        Simple image viewer
      </div>}
    </>
  )
}