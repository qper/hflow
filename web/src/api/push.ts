import { apiClient } from './client';

export async function savePushSubscription(subscription: PushSubscription) {
  const response = await apiClient.POST('/api/v1/push/subscriptions', {
    body: {
      endpoint: subscription.endpoint,
      keys: {
        p256dh: btoa(String.fromCharCode(...new Uint8Array(subscription.getKey('p256dh')!))),
        auth: btoa(String.fromCharCode(...new Uint8Array(subscription.getKey('auth')!))),
      },
    },
  });

  if (response.error) {
    throw new Error('Unable to save push subscription');
  }
}
