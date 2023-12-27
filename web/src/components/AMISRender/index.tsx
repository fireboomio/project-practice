import { useLocation, useParams, useSearchParams } from '@modern-js/runtime/router'
import { render as amisRender } from 'amis'

import { StorageKey, tokenStorage, alert, copy, fetcher, notify, theme } from '../../utils'
import { useAdmin } from '@/store/admin'

interface AMISRenderProps {
  schema: any
  transform?: Function
  session?: string
  [key: string]: any
}

const AMISRender = ({ schema, transform, session, ...rest }: AMISRenderProps) => {
  const adminStore = useAdmin()
  const params = useParams()
  const { pathname } = useLocation()
  const [searchParams] = useSearchParams()
  const _schema = schema.type ? schema : { type: 'page', ...schema }

  return (
    <div>
      {amisRender(
        _schema,
        {
          data: {
            pathname,
            query: Object.fromEntries(searchParams),
            params,
            TOKEN: tokenStorage.getItem(StorageKey.AccessToken),
            BASE_URL: '',
            adminStore
          },
          propsTransform: transform,
          ...rest
        },
        {
          fetcher,
          isCancel: (value: any) => !!value.__CANCEL__,
          notify,
          alert,
          copy,
          theme
        }
      )}
    </div>
  )
}

export default AMISRender
