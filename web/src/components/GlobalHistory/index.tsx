import { useNavigate, NavigateFunction } from 'react-router-dom'

export let globalNavigate: NavigateFunction

export const GlobalHistory = () => {
  globalNavigate = useNavigate()

  return null
};