import type { ProfileStatusValue, RequirementStatusValue } from '@/types/api'

export function statusToBadge(status: RequirementStatusValue | ProfileStatusValue) {
  if (status === 'Valid' || status === 'Eligible') return 'success'
  if (status === 'Expiring' || status === 'At Risk') return 'warning'
  if (status === 'Expired' || status === 'Not Eligible') return 'destructive'
  return 'secondary'
}
