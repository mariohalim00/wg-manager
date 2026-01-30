<script lang="ts">
  export let text: string;
  // In a real app, you would use a library to generate the QR code SVG or canvas
  const qrCodeUrl = `https://api.qrserver.com/v1/create-qr-code/?data=${encodeURIComponent(text)}&size=200x200`;
</script>

<img src={qrCodeUrl} alt="QR Code for {text}" />
