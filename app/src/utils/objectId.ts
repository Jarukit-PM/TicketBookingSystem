const ZERO_OBJECT_ID = '000000000000000000000000'

/** True when id is a non-empty MongoDB ObjectId that is not the zero value. */
export function isValidObjectId(id?: string | null): boolean {
  return Boolean(id && id !== ZERO_OBJECT_ID)
}
