import { ProjectDetail } from "@/components/project/project-detail"

interface ProjectPageProps {
  params: Promise<{
    id: string
  }>
}

export default async function ProjectPage({ params }: ProjectPageProps) {
  const { id } = await params
  
  return (
    <div className="min-h-screen bg-background sandstorm-bg">
      <ProjectDetail projectId={id} />
    </div>
  )
}
