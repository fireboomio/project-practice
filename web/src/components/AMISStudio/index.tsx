import './index.scss'

// import 'amis/lib/themes/antd.css'
// import 'amis/lib/helper.css'
// import 'amis/sdk/iconfont.css'
// import 'amis-editor-core/lib/style.css'
import type { SchemaObject } from 'amis'
import { Editor, ShortcutKey } from 'amis-editor'
import { useState } from 'react'

// import { ReactComponent as IconH5Preview } from '@/assets/h5-preview.svg'
// import { ReactComponent as IconPCPreview } from '@/assets/pc-preview.svg'
import { useAdmin } from '@/store/admin'
import { alert, copy, fetcher, notify, theme } from '@/utils/amis'

const host = `${window.location.protocol}//${window.location.host}`
const schemaUrl = `${host}/#/schema.json`

export type AMISStudioProps = {
  schema?: string | SchemaObject
  onSave: (schema?: SchemaObject) => Promise<void>
  onExit: () => void
}

const AMISStudio = ({ schema: originSchema, onSave, onExit }: AMISStudioProps) => {
  const [schema, setSchema] = useState<SchemaObject>(
    typeof originSchema === 'string' ? JSON.parse(originSchema) : originSchema
  )
  const { isMobile, setIsMobile, preview, setPreview } = useAdmin()

  const exit = () => {
    if (preview) {
      setPreview(false)
    }
    onExit()
  }

  const save = async () => {
    if (schema) {
      schema.type = 'page'
      await onSave(schema)
    }
  }

  const saveAndExit = async () => {
    await save()
    exit()
  }

  return (
    <div className="amis-studio">
      <div className="studio-header">
        <div className="studio-title">amis 可视化编辑器</div>
        <div className="studio-view-mode-group-container">
          <div className="studio-view-mode-group">
            {/* <div
              className={`studio-view-mode-btn editor-header-icon ${!isMobile ? 'is-active' : ''}`}
              onClick={() => {
                setIsMobile(false)
              }}
            >
              <IconH5Preview />
            </div>
            <div
              className={`studio-view-mode-btn editor-header-icon ${isMobile ? 'is-active' : ''}`}
              onClick={() => {
                setIsMobile(true)
              }}
            >
              <IconPCPreview />
            </div> */}
          </div>
        </div>

        <div className="studio-header-actions">
          <ShortcutKey />
          <div
            className={`header-action-btn m-1 ${preview ? 'primary' : ''}`}
            onClick={() => {
              setPreview(!preview)
            }}
          >
            {preview ? '编辑' : '预览'}
          </div>
          {!preview && (
            <>
              <div className={`header-action-btn exit-btn`} onClick={exit}>
                退出不保存
              </div>
              <div className={`header-action-btn`} onClick={saveAndExit}>
                保存并退出
              </div>
            </>
          )}
        </div>
      </div>
      <div className="studio-inner">
        <Editor
          theme={theme}
          preview={preview}
          isMobile={isMobile}
          // value={schema}
          onChange={s => {
            setSchema(s)
          }}
          onPreview={() => {
            setPreview(true)
          }}
          onSave={save}
          className="is-fixed"
          $schemaUrl={schemaUrl}
          showCustomRenderersPanel={true}
          amisEnv={{
            // @ts-ignore
            fetcher,
            notify,
            alert,
            copy
          }}
        />
      </div>
    </div>
  )
}

export default AMISStudio
