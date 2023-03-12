type Props = {
  hideSubText?: boolean
}

export default function Logo(props: Props) {
  return (
    <>
      <h1 class='pt-7 text-center font-extrabold text-transparent text-6xl bg-clip-text bg-gradient-to-r from-blue-600 to-sky-500'>
        Fuu
      </h1>
      {!props.hideSubText && <div class='text-center text-xl text-neutral-100'>
        Simple image viewer
      </div>}
    </>
  )
}